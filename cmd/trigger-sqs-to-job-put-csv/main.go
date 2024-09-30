package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/Mad-Pixels/lingocards-api/pkg/trigger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")

	s3Bucket *cloud.Bucket
	dbDynamo *cloud.Dynamo
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

const filenameKey = "filename"

func handler(ctx context.Context, log zerolog.Logger, record any) error {
	dynamoRecord, ok := record.(events.DynamoDBEventRecord)
	if !ok {
		return fmt.Errorf("invalid record type: %T", record)
	}

	if dynamoRecord.Change.NewImage != nil {
		key, ok := dynamoRecord.Change.NewImage[filenameKey]
		if !ok {
			log.Error().Msg("Attribute 'filename' not found in new image")
			return fmt.Errorf("attribute 'filename' not found in new image")
		}

		fileBody, err := s3Bucket.Get(ctx, key.String(), serviceProcessingBucket)
		if err != nil {
			log.Error().Err(err).
				Str("bucket", serviceProcessingBucket).
				Str("key", key.String()).
				Msg("Failed to get file from bucket")
			return err
		}
		defer fileBody.Close()

		err = s3Bucket.Put(ctx, key.String(), serviceDictionaryBucket, fileBody.(io.Reader), "application/octet-stream")
		if err != nil {
			log.Error().Err(err).
				Str("bucket", serviceDictionaryBucket).
				Str("key", key.String()).
				Msg("Failed to upload file to bucket")
			return err
		}

		log.Info().
			Str("file", key.String()).
			Str("from", serviceProcessingBucket).
			Str("to", serviceDictionaryBucket).
			Msg("Successfully moved file")
	}

	return nil
}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{MaxWorkers: 4},
			handler,
		).Handle,
	)
}
