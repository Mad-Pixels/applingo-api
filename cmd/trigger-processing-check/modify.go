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

	// actions
	needProcess bool

	// data
	oldItem *applingoprocessing.SchemaItem
	newItem *applingoprocessing.SchemaItem
}

func detectChanges(old, newItem *applingoprocessing.SchemaItem) changes {
	c := changes{
		score:  old.Score != newItem.Score,
		upload: old.Upload != newItem.Upload,

		oldItem: old,
		newItem: newItem,
	}

	if c.upload && newItem.Upload == 1 {
		c.needProcess = true
	}
	if c.score && newItem.Score >= autoUploadScoreThreshold {
		c.needProcess = true
	}
	return c
}

func modify(ctx context.Context, e events.DynamoDBEventRecord) error {
	newItem, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(e)
	if err != nil {
		return fmt.Errorf("failed extract new data: %w", err)
	}
	oldItem, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(events.DynamoDBEventRecord{
		Change: events.DynamoDBStreamRecord{
			NewImage: e.Change.OldImage,
		},
	})
	if err != nil {
		return fmt.Errorf("failed extract old data: %w", err)
	}
	c := detectChanges(oldItem, newItem)

	if c.needProcess {
		if err := processRecordToDictionary(ctx, &c); err != nil {
			return fmt.Errorf("failed to process record: %w", err)
		}
		if err := updateRecordUploadStatus(ctx, &c); err != nil {
			return fmt.Errorf("failed to update record upload status: %w", err)
		}
	}
	return nil
}

func updateRecordUploadStatus(ctx context.Context, c *changes) error {
	key, err := applingoprocessing.CreateKeyFromItem(*c.newItem)
	if err != nil {
		return fmt.Errorf("failed to create key for item: %w", err)
	}
	update := expression.
		Set(
			expression.Name(applingoprocessing.ColumnUpload),
			expression.Value(applingoprocessing.BoolToInt(true)),
		)
	return dbDynamo.Update(
		ctx,
		applingoprocessing.TableName,
		key,
		update,
		expression.AttributeExists(expression.Name(applingoprocessing.ColumnId)),
	)
}

func processRecordToDictionary(ctx context.Context, c *changes) error {
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
	dictionaryFileID := utils.RecordToFileID(c.newItem.Id)

	// copy dictionary data.
	if err = s3Bucket.Copy(ctx, dictionaryFileID, serviceProcessingBucket, dictionaryFileID, serviceDictionaryBucket); err != nil {
		return fmt.Errorf("failed to copy dictionary from processing to service: %w", err)
	}
	if err = s3Bucket.WaitOrError(ctx, dictionaryFileID, serviceDictionaryBucket, 3, 200*time.Millisecond); err != nil {
		return fmt.Errorf("failed to check object in bucket")
	}

	// insert data to dynamoDB.
	if err := dbDynamo.Put(
		ctx,
		applingodictionary.TableSchema.TableName,
		dynamoItem,
		expression.AttributeNotExists(expression.Name(applingodictionary.ColumnId)),
	); err != nil {
		s3Err := s3Bucket.Delete(ctx, dictionaryFileID, serviceProcessingBucket)
		if s3Err != nil {
			return fmt.Errorf("failed add new dictionary in dynamoDB: %w, also cannot delete dictionary from bucket: %w", err, s3Err)
		}
		return fmt.Errorf("failed add new dictionary in dynamoDB: %w, dictionary was removed from bucket", err)
	}
	return nil
}
