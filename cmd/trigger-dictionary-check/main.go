package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/forge"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/rs/zerolog"
)

const (
	lambdaWatchdog           = 180 * time.Second
	defaultBackoff           = 15 * time.Second
	autoUploadScoreThreshold = 90
	defaultRetries           = 2
	defaultMaxWorkers        = 5
)

var (
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceForgeBucket      = os.Getenv("SERVICE_FORGE_BUCKET")
	lambdaTimeout           = os.Getenv("LAMBDA_TIMEOUT_SECONDS")
	awsRegion               = os.Getenv("AWS_REGION")
	openaiToken             = os.Getenv("OPENAI_KEY")

	gptClient *chatgpt.Client
	dbDynamo  *cloud.Dynamo
	s3Bucket  *cloud.Bucket

	timeout = utils.GetTimeout(lambdaTimeout, lambdaWatchdog)
)

func init() {
	debug.SetGCPercent(500)

	gptClient = chatgpt.MustClient(
		httpclient.New().
			WithTimeout(timeout).
			WithMaxRetries(defaultRetries, defaultBackoff).
			WithRetryCondition(func(statusCode int, responseBody string) bool {
				return statusCode >= 500 && statusCode < 600
			}),
		openaiToken,
	)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
	dbDynamo = cloud.NewDynamo(cfg)
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	var dynamoDBEvent events.DynamoDBEventRecord
	if err := serializer.UnmarshalJSON(record, &dynamoDBEvent); err != nil {
		return fmt.Errorf("failed to unmarshal request record: %w", err)
	}

	switch dynamoDBEvent.EventName {
	case "INSERT":
		item, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(dynamoDBEvent)
		if err != nil {
			return fmt.Errorf("failed to extract item from DynamoDB event: %w", err)
		}

		req := forge.NewRequestDictionaryCheck()
		result, err := forge.Check(ctx, req, item, serviceForgeBucket, serviceProcessingBucket, gptClient, s3Bucket)
		if err != nil {
			log.Error().Err(err).Str("file", item.File).Msg("failed to check dictionary")
			return fmt.Errorf("failed to check dictionary: %w", err)
		}

		key, err := applingoprocessing.CreateKeyFromItem(*item)
		if err != nil {
			log.Error().Err(err).Str("id", item.Id).Msg("failed to create key for item")
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
		if err = dbDynamo.Update(ctx, applingoprocessing.TableSchema.TableName, key, update, condition); err != nil {
			log.Error().Err(err).Str("id", item.Id).Msg("failed to update item in DynamoDB")
			return fmt.Errorf("failed to update item in DynamoDB: %w", err)
		}

	case "MODIFY":
		item, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(dynamoDBEvent)
		if err != nil {
			return fmt.Errorf("failed to extract item from DynamoDB event: %w", err)
		}
		if item.Score >= autoUploadScoreThreshold || applingoprocessing.IntToBool(item.Upload) {
			shemaItem := applingodictionary.SchemaItem{
				Id:          item.Id,
				Subcategory: item.Subcategory,

				IsPublic:  applingodictionary.BoolToInt(true),
				Created:   int(time.Now().Unix()),
				Category:  "Languages",
				Downloads: 0,
				Rating:    0,

				Description: item.Overview,
				Author:      item.Author,
				Name:        item.Name,
				Topic:       item.Topic,
				Level:       item.Level,
				Words:       item.Words,
				Filename:    fmt.Sprintf("%s.json", item.File),
				Dictionary:  item.File,

				LevelSubcategoryIsPublic: fmt.Sprintf("%s#%s#%d", item.Level, item.Subcategory, applingodictionary.BoolToInt(true)),
				SubcategoryIsPublic:      fmt.Sprintf("%s#%d", item.Subcategory, applingodictionary.BoolToInt(true)),
				LevelIsPublic:            fmt.Sprintf("%s#%d", item.Level, applingodictionary.BoolToInt(true)),
			}
			dynamoItem, err := applingodictionary.PutItem(shemaItem)
			if err != nil {
				return fmt.Errorf("failed prepare dynamo item: %w", err)
			}

			if err := s3Bucket.Move(ctx, item.File, serviceProcessingBucket, fmt.Sprintf("%s.json", item.File), serviceDictionaryBucket); err != nil {
				return fmt.Errorf("failed to move dictionary from processing to service: %w", err)
			}
			if err := dbDynamo.Put(
				ctx,
				applingodictionary.TableSchema.TableName,
				dynamoItem,
				expression.AttributeNotExists(expression.Name(applingodictionary.ColumnId)),
			); err != nil {
				s3Err := s3Bucket.Delete(ctx, item.File, serviceProcessingBucket)
				if s3Err != nil {
					return fmt.Errorf("failed add new dictionary in dynamoDB: %w, also cannot delete dictionary from bucket: %w", err, s3Err)
				}
				return fmt.Errorf("failed add new dictionary in dyynamoDB: %w, dictionary was removed from bucket", err)
			}
		}
		log.Info().Str("id", item.Id).Msg("item modified, no action required")

	case "REMOVE":
		var fileId string
		if fileKey, ok := dynamoDBEvent.Change.Keys["file"]; ok {
			fileId = fileKey.String()
		}
		if fileId == "" {
			log.Warn().Msg("file key is empty in REMOVE event, cannot delete file")
			return nil
		}

		if err := s3Bucket.Delete(ctx, fileId, serviceProcessingBucket); err != nil {
			log.Error().Err(err).Str("file", fileId).Msg("failed to delete file from bucket")
			return fmt.Errorf("failed to delete file from bucket: %w", err)
		}
		log.Info().Str("file", fileId).Msg("file deleted successfully")
		return nil

	default:
		log.Warn().Str("eventName", dynamoDBEvent.EventName).Msg("unhandled event type")
	}
	return nil
}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{MaxWorkers: 6},
			handler,
		).Handle,
	)
}
