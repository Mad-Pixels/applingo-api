package main

import (
	"context"
	"encoding/json"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "failed to unmarshal DynamoDB event")
	}
	item, err := applingoprocessing.ExtractFromDynamoDBStreamEvent(dynamoDBEvent)
	if err != nil {
		return errors.Wrap(err, "failed to extract item from DynamoDB event")
	}

	switch dynamoDBEvent.EventName {
	case "INSERT":
		log.Error().Any("item", item).Msg("INSERT")
	case "MODIFY":
		log.Error().Any("item", item).Msg("MODIFY")
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
