package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

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
	defaultBackoff    = 15 * time.Second
	defaultRetries    = 2
	defaultMaxWorkers = 5
)

var (
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	serviceForgeBucket      = os.Getenv("SERVICE_FORGE_BUCKET")
	lambdaTimeout           = os.Getenv("LAMBDA_TIMEOUT_SECONDS")
	awsRegion               = os.Getenv("AWS_REGION")
	openaiToken             = os.Getenv("OPENAI_KEY")

	gptClient *chatgpt.Client
	dbDynamo  *cloud.Dynamo
	s3Bucket  *cloud.Bucket

	timeout = utils.GetTimeout(lambdaTimeout, 240*time.Second)
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
	item, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(dynamoDBEvent)
	if err != nil {
		return fmt.Errorf("failed to extract item from DynamoDB event: %w", err)
	}

	switch dynamoDBEvent.EventName {
	case "INSERT":
		checkReq := forge.NewRequestDictionaryCheck()
		checkReq.DictionaryFile = item.File

		result, err := forge.Check(ctx, checkReq, serviceForgeBucket, serviceProcessingBucket, gptClient, s3Bucket)
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
				expression.Value(result.Meta.Score),
			).
			Set(
				expression.Name(applingoprocessing.ColumnReason),
				expression.Value(result.Meta.Reason),
			)

		condition := expression.AttributeExists(expression.Name(applingoprocessing.ColumnId))
		if err = dbDynamo.Update(ctx, applingoprocessing.TableSchema.TableName, key, update, condition); err != nil {
			log.Error().Err(err).Str("id", item.Id).Msg("failed to update item in DynamoDB")
			return fmt.Errorf("failed to update item in DynamoDB: %w", err)
		}

	case "MODIFY":
		log.Info().Str("id", item.Id).Msg("item modified, no action required")

	case "REMOVE":
		log.Info().Str("id", item.Id).Msg("item removed, no action required")

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
