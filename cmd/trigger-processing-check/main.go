// Package main provides a Lambda function that handles processing of dictionary records,
// including insertion, modification, and deletion. It integrates with AWS DynamoDB and S3
// and uses OpenAI for validating dictionary content.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
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
			WithRetryCondition(func(statusCode int, _ string) bool {
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
		log.Info().Msg("INSERT event")
		if err := insert(ctx, dynamoDBEvent); err != nil {
			return fmt.Errorf("failed process item: %w", err)
		}
	case "MODIFY":
		log.Info().Msg("MODIFY event")
		if err := modify(ctx, dynamoDBEvent); err != nil {
			return fmt.Errorf("failed modify item: %w", err)
		}
	case "REMOVE":
		log.Info().Msg("REMOVE event")
		if err := remove(ctx, dynamoDBEvent); err != nil {
			return fmt.Errorf("failed to delete file from bucket: %w", err)
		}
	default:
	}
	return nil
}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{MaxWorkers: defaultMaxWorkers},
			handler,
		).Handle,
	)
}
