package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/persistent"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type changes struct {
	// check
	description bool
	overview    bool
	topic       bool
	level       bool
	languages   bool
	words       bool
	subcategory bool
	promptCraft bool
	promptCheck bool

	// process
	score  bool
	upload bool

	// action flags.
	needProcess bool
	needUpdate  bool

	// data
	oldItem *applingoprocessing.SchemaItem
	newItem *applingoprocessing.SchemaItem
}

func detectChanges(old, new *applingoprocessing.SchemaItem) changes {
	c := changes{
		score:  old.Score != new.Score,
		upload: old.Upload != new.Upload,

		description: old.Description != new.Description,
		overview:    old.Overview != new.Overview,
		topic:       old.Topic != new.Topic,
		level:       old.Level != new.Level,
		languages:   old.Languages != new.Languages,
		words:       old.Words != new.Words,
		subcategory: old.Subcategory != new.Subcategory,
		promptCraft: old.PromptCraft != new.PromptCraft,
		promptCheck: old.PromptCheck != new.PromptCheck,

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
	check(&c)

	if c.needProcess {
		if err := processRecord(ctx, &c); err != nil {
			return fmt.Errorf("failed to process record: %w", err)
		}
	}
	if c.needUpdate {
		if err := updateRecord(ctx, *c.newItem); err != nil {
			return fmt.Errorf("failed to update record with validated data: %w", err)
		}
	}
	return nil
}

func check(c *changes) {
	beforeChecks := *c.newItem

	if c.description {
		c.newItem.Description = persistent.FieldDescriptionOrDefault(c.newItem.Description, c.oldItem.Description)
	}
	if c.topic {
		c.newItem.Topic = persistent.FieldTopicOrDefault(c.newItem.Topic, c.oldItem.Topic)
	}
	if c.level {
		c.newItem.Level = persistent.FieldLevelOrDefault(c.newItem.Level, c.oldItem.Level)
	}
	if c.languages {
		c.newItem.Languages = persistent.FieldLanguagesOrDefault(c.newItem.Languages, c.oldItem.Languages)
	}
	if c.words {
		c.newItem.Words = persistent.FieldWordsOrDefault(c.newItem.Words, c.oldItem.Words)
	}
	if c.subcategory {
		c.newItem.Subcategory = persistent.FieldSubcategoryOrDefault(c.newItem.Subcategory, c.oldItem.Subcategory)
	}
	if c.promptCraft {
		c.newItem.PromptCraft = persistent.FieldPromtOrDefault(c.newItem.PromptCraft, c.oldItem.PromptCraft)
	}
	if c.promptCheck {
		c.newItem.PromptCheck = persistent.FieldPromtOrDefault(c.newItem.PromptCheck, c.oldItem.PromptCheck)
	}
	c.needUpdate = !reflect.DeepEqual(beforeChecks, *c.newItem)
}

func updateRecord(ctx context.Context, item applingoprocessing.SchemaItem) error {
	condition := expression.And(
		expression.AttributeExists(expression.Name(applingoprocessing.ColumnId)),
		expression.AttributeExists(expression.Name(applingoprocessing.ColumnCreated)),
	)
	dynamoItem, err := applingoprocessing.PutItem(item)
	if err != nil {
		return err
	}
	return dbDynamo.Put(
		ctx,
		applingoprocessing.TableSchema.TableName,
		dynamoItem,
		condition,
	)
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
