package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"
	"github.com/Mad-Pixels/lingocards-api/pkg/trigger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	dictionaryFilenameKey = "dictionary"
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

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	var sqsRecord events.SQSMessage
	if err := serializer.UnmarshalJSON(record, &sqsRecord); err != nil {
		return errors.Wrap(err, "failed to unmarshal SQS record")
	}
	var dynamoDBEvent events.DynamoDBEventRecord
	if err := serializer.UnmarshalJSON([]byte(sqsRecord.Body), &dynamoDBEvent); err != nil {
		return errors.Wrap(err, "failed to unmarshal DynamoDB event from SQS message body")
	}

	bucketKey, ok := dynamoDBEvent.Change.NewImage[dictionaryFilenameKey]
	if !ok {
		return errors.New("'dictionaryFilenameKey' not found in DynamoDB event")
	}
	if bucketKey.DataType() != events.DataTypeString {
		return errors.New("'dictionaryFilenameKey' is not a string in DynamoDB event")
	}
	return processFile(ctx, bucketKey.String())
}

func processFile(ctx context.Context, filename string) error {
	fileBody, err := s3Bucket.Get(ctx, filename, serviceProcessingBucket)
	if err != nil {
		return errors.Wrapf(err, "failed to get file %s from bucket %s", filename, serviceProcessingBucket)
	}
	defer fileBody.Close()

	if err = s3Bucket.Put(ctx, filename, serviceDictionaryBucket, fileBody.(io.Reader), "application/octet-stream"); err != nil {
		return errors.Wrapf(err, "failed to upload file %s to bucket %s", filename, serviceDictionaryBucket)
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
