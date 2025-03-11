package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type changes struct {
	// process
	score  bool
	upload bool

	// general flags.
	needProcess bool

	// data
	oldItem *applingoprocessing.SchemaItem
	newItem *applingoprocessing.SchemaItem
}

func detectChanges(old, new *applingoprocessing.SchemaItem) changes {
	c := changes{
		score:  old.Score != new.Score,
		upload: old.Upload != new.Upload,

		oldItem: old,
		newItem: new,
	}

	if c.upload && new.Upload == 1 {
		c.needProcess = true
	}
	if c.score && new.Score >= autoUploadScoreThreshold {
		c.needProcess = true
	}
	return c
}

func modify(ctx context.Context, e events.DynamoDBEventRecord) error {
	new, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(e)
	if err != nil {
		return fmt.Errorf("failed extract new data: %w", err)
	}
	old, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(events.DynamoDBEventRecord{
		Change: events.DynamoDBStreamRecord{
			NewImage: e.Change.OldImage,
		},
	})
	if err != nil {
		return fmt.Errorf("failed extract old data: %w", err)
	}
	c := detectChanges(old, new)

	if c.needProcess {
		if err := processRecord(ctx, &c); err != nil {
			return fmt.Errorf("failed to process record: %w", err)
		}
	}
	return nil
}

func processRecord(ctx context.Context, c *changes) error {
	existingKey, err := applingodictionary.CreateKey(c.newItem.Id, c.newItem.Subcategory)
	if err != nil {
		return fmt.Errorf("failed to create key for checking: %w", err)
	}
	exists, err := dbDynamo.Exists(ctx, applingodictionary.TableSchema.TableName, existingKey)
	if err != nil {
		return fmt.Errorf("failed to check if record exists: %w", err)
	}
	if exists {
		fmt.Printf("Record with ID %s already exists in target table, skipping\n", c.newItem.Id)
		return nil
	}

	// prepare dynamo item.
	schemaItem := applingodictionary.SchemaItem{
		// identifier.
		Id:          c.newItem.Id,
		Subcategory: c.newItem.Subcategory,

		// data fields.
		Description: c.newItem.Overview,
		Author:      c.newItem.Author,
		Name:        c.newItem.Name,
		Topic:       c.newItem.Topic,
		Level:       c.newItem.Level,
		Words:       c.newItem.Words,

		// system fileds.
		IsPublic: applingodictionary.BoolToInt(true),
		Created:  int(time.Now().Unix()),
		Category: "Languages",

		// composite keys.
		LevelSubcategoryIsPublic: fmt.Sprintf("%s#%s#%d", c.newItem.Level, c.newItem.Subcategory, applingodictionary.BoolToInt(true)),
		SubcategoryIsPublic:      fmt.Sprintf("%s#%d", c.newItem.Subcategory, applingodictionary.BoolToInt(true)),
		LevelIsPublic:            fmt.Sprintf("%s#%d", c.newItem.Level, applingodictionary.BoolToInt(true)),
	}
	dynamoItem, err := applingodictionary.PutItem(schemaItem)
	if err != nil {
		return fmt.Errorf("failed prepare dynamo item: %w", err)
	}

	// copy dictionary data.
	if err = s3Bucket.Copy(ctx, c.newItem.Id, serviceProcessingBucket, utils.RecordToFileID(c.newItem.Id), serviceDictionaryBucket); err != nil {
		return fmt.Errorf("failed to copy dictionary from processing to service: %w", err)
	}
	if err = s3Bucket.WaitOrError(ctx, utils.RecordToFileID(c.newItem.Id), serviceDictionaryBucket, 3, 200*time.Millisecond); err != nil {
		return fmt.Errorf("failed to check object in bucket")
	}

	// insert data to dynamoDB.
	if err := dbDynamo.Put(
		ctx,
		applingodictionary.TableSchema.TableName,
		dynamoItem,
		expression.AttributeNotExists(expression.Name(applingodictionary.ColumnId)),
	); err != nil {
		s3Err := s3Bucket.Delete(ctx, utils.RecordToFileID(c.newItem.Id), serviceProcessingBucket)
		if s3Err != nil {
			return fmt.Errorf("failed add new dictionary in dynamoDB: %w, also cannot delete dictionary from bucket: %w", err, s3Err)
		}
		return fmt.Errorf("failed add new dictionary in dynamoDB: %w, dictionary was removed from bucket", err)
	}
	return nil
}
