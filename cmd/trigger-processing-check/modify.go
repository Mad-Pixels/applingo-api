package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type changes struct {
	id          bool
	filename    bool
	name        bool
	score       bool
	upload      bool
	level       bool
	languages   bool
	words       bool
	subcategory bool
	promptCraft bool
	promptCheck bool
	description bool
	topic       bool
	overview    bool

	changed bool
	process bool
	oldItem *applingoprocessing.SchemaItem
	newItem *applingoprocessing.SchemaItem
}

func detectChanges(oldItem, newItem *applingoprocessing.SchemaItem) changes {
	c := changes{
		oldItem: oldItem,
		newItem: newItem,
	}

	c.subcategory = oldItem.Subcategory != newItem.Subcategory
	c.promptCraft = oldItem.PromptCraft != newItem.PromptCraft
	c.promptCheck = oldItem.PromptCheck != newItem.PromptCheck
	c.description = oldItem.Description != newItem.Description
	c.languages = oldItem.Languages != newItem.Languages
	c.overview = oldItem.Overview != newItem.Overview
	c.filename = oldItem.Filename != newItem.Filename
	c.upload = oldItem.Upload != newItem.Upload
	c.score = oldItem.Score != newItem.Score
	c.level = oldItem.Level != newItem.Level
	c.words = oldItem.Words != newItem.Words
	c.topic = oldItem.Topic != newItem.Topic
	c.name = oldItem.Name != newItem.Name
	c.id = oldItem.Id != newItem.Id

	c.changed = c.subcategory ||
		c.promptCraft ||
		c.description ||
		c.languages ||
		c.filename ||
		c.level ||
		c.words ||
		c.topic ||
		c.name ||
		c.id

	c.process = c.score || c.upload
	return c
}

func modify(ctx context.Context, e events.DynamoDBEventRecord) error {
	old, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(events.DynamoDBEventRecord{
		Change: events.DynamoDBStreamRecord{
			NewImage: e.Change.OldImage,
		},
	})
	fmt.Println("Event Name:", e.EventName)
	fmt.Println("Has OldImage:", e.Change.OldImage != nil)
	fmt.Println("Has NewImage:", e.Change.NewImage != nil)

	if e.Change.OldImage == nil {
		return fmt.Errorf("OldImage is nil for MODIFY event")
	}
	if err != nil {
		return nil
	}

	new, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(e)
	if err != nil {
		return fmt.Errorf("failed extract new data: %w", err)
	}

	fmt.Printf("BEFORE DETECT - Old: %+v\n", old)
	fmt.Printf("BEFORE DETECT - New: %+v\n", new)

	c := detectChanges(old, new)
	fmt.Printf("AFTER DETECT - Changed fields: score=%v, upload=%v, description=%v, ...\n",
		c.score, c.upload, c.description)
	fmt.Printf("AFTER DETECT - Changed overall: %v, Process: %v\n", c.changed, c.process)
	if c.changed {
		err = invokeChanges(ctx, &c)
		if err != nil {
			return err
		}
	}
	if c.process {
		err = invokeProcess(ctx, &c)
	}
	return err
}

func invokeChanges(ctx context.Context, c *changes) error {
	if c == nil {
		return errors.New("incomming changes is nil")
	}
	key, err := applingoprocessing.CreateKeyFromItem(*c.oldItem)
	if err != nil {
		return fmt.Errorf("failed to create key for item: %w", err)
	}

	// Инициализируем пустой UpdateBuilder
	var builder expression.UpdateBuilder
	hasUpdate := false

	if c.description {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnDescription),
			expression.Value(c.oldItem.Description),
		)
		hasUpdate = true
	}
	if c.promptCraft {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnPromptCraft),
			expression.Value(c.oldItem.PromptCraft),
		)
		hasUpdate = true
	}
	if c.languages {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnLanguages),
			expression.Value(c.oldItem.Languages),
		)
		hasUpdate = true
	}
	if c.filename {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnFilename),
			expression.Value(utils.RecordToFileID(utils.GenerateDictionaryID(c.newItem.Name, c.newItem.Author))),
		)
		hasUpdate = true
	}
	if c.level {
		level := ""

		lvl, err := types.ParseLanguageLevel(strings.ToUpper(c.newItem.Level))
		if err != nil {
			level = c.oldItem.Level
		} else {
			level = lvl.String()
		}
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnLevel),
			expression.Value(level),
		)
		hasUpdate = true
	}
	if c.words {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnWords),
			expression.Value(c.oldItem.Words),
		)
		hasUpdate = true
	}
	if c.topic {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnTopic),
			expression.Value(c.oldItem.Topic),
		)
		hasUpdate = true
	}
	if c.name {
		builder = builder.Set(
			expression.Name(applingoprocessing.ColumnName),
			expression.Value(c.newItem.Name),
		)
		hasUpdate = true
	}

	// Если нет обновлений, просто выходим
	if !hasUpdate {
		return nil
	}

	condition := expression.AttributeExists(expression.Name(applingoprocessing.ColumnId))
	return dbDynamo.Update(ctx, applingoprocessing.TableSchema.TableName, key, builder, condition)
}

func invokeProcess(ctx context.Context, c *changes) error {
	if c == nil {
		return errors.New("incomming changes is nil")
	}

	needProcess := false
	if c.oldItem.Score != c.newItem.Score && c.newItem.Score >= autoUploadScoreThreshold {
		needProcess = true
	}
	if c.oldItem.Upload != c.newItem.Upload && c.newItem.Upload == 1 {
		needProcess = true
	}

	if !needProcess {
		return nil
	}
	shemaItem := applingodictionary.SchemaItem{
		Id:          c.newItem.Id,
		Subcategory: c.newItem.Subcategory,

		Description: c.newItem.Overview,
		Dictionary:  utils.RecordToFileID(c.newItem.Id),
		Author:      c.newItem.Author,
		Name:        c.newItem.Name,
		Topic:       c.newItem.Topic,
		Level:       c.newItem.Level,
		Words:       c.newItem.Words,

		IsPublic: applingodictionary.BoolToInt(true),
		Created:  int(time.Now().Unix()),
		Category: "Languages",

		LevelSubcategoryIsPublic: fmt.Sprintf("%s#%s#%d", c.newItem.Level, c.newItem.Subcategory, applingodictionary.BoolToInt(true)),
		SubcategoryIsPublic:      fmt.Sprintf("%s#%d", c.newItem.Subcategory, applingodictionary.BoolToInt(true)),
		LevelIsPublic:            fmt.Sprintf("%s#%d", c.newItem.Level, applingodictionary.BoolToInt(true)),
	}
	dynamoItem, err := applingodictionary.PutItem(shemaItem)
	if err != nil {
		return fmt.Errorf("failed prepare dynamo item: %w", err)
	}

	if err = s3Bucket.Copy(ctx, c.newItem.Id, serviceProcessingBucket, utils.RecordToFileID(c.newItem.Id), serviceDictionaryBucket); err != nil {
		return fmt.Errorf("failed to copy dictionary from processing to service: %w", err)
	}
	if err = s3Bucket.WaitOrError(ctx, utils.RecordToFileID(c.newItem.Id), serviceDictionaryBucket, 3, 200*time.Millisecond); err != nil {
		return fmt.Errorf("failed to check object in bucket")
	}

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
		return fmt.Errorf("failed add new dictionary in dyynamoDB: %w, dictionary was removed from bucket", err)
	}
	return nil
}
