package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
)

const (
	defaultMaxWorkers = 5
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")

	dbDynamo *cloud.Dynamo
	s3Bucket *cloud.Bucket
)

func init() {
	debug.SetGCPercent(500)

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

	case "MODIFY":

	case "REMOVE":
		var fileId string
		fmt.Println(dynamoDBEvent)
		fmt.Println(dynamoDBEvent.Change.Keys)
		if fileKey, ok := dynamoDBEvent.Change.Keys["id"]; ok {
			fileId = fileKey.String() + ".json"
		}
		if fileId == "" {
			log.Warn().Msg("file key is empty in REMOVE event, cannot delete file")
			return nil
		}

		if err := s3Bucket.Delete(ctx, fileId, serviceDictionaryBucket); err != nil {
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
			trigger.Config{MaxWorkers: defaultMaxWorkers},
			handler,
		).Handle,
	)
}
