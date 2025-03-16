package main

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/forge"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func insert(ctx context.Context, e events.DynamoDBEventRecord) error {
	item, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(e)
	if err != nil {
		return fmt.Errorf("failed to extract item from DynamoDB event: %w", err)
	}

	req := forge.NewRequestDictionaryCheck()
	result, err := forge.Check(ctx, req, item, serviceForgeBucket, serviceProcessingBucket, gptClient, s3Bucket)
	if err != nil {
		return fmt.Errorf("failed to check dictionary: %w", err)
	}

	key, err := applingoprocessing.CreateKeyFromItem(*item)
	if err != nil {
		return fmt.Errorf("failed to create key for item: %w", err)
	}
	update := expression.
		Set(
			expression.Name(applingoprocessing.ColumnScore),
			expression.Value(result.GetScore()),
		).
		Set(
			expression.Name(applingoprocessing.ColumnReason),
			expression.Value(result.GetReason()),
		).
		Set(
			expression.Name(applingoprocessing.ColumnPromptCheck),
			expression.Value(utils.JoinValues(result.GetPrompt(), string(result.GetModel()))),
		)
	condition := expression.AttributeExists(expression.Name(applingoprocessing.ColumnId))
	return dbDynamo.Update(ctx, applingoprocessing.TableSchema.TableName, key, update, condition)
}
