package main

import (
	"context"
	"encoding/json"
	"fmt"
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

const filenameKey = "filename"

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
	var dynamoRecord events.DynamoDBEventRecord
	if err := serializer.UnmarshalJSON(record, &dynamoRecord); err != nil {
		return errors.Wrap(err, "failed to unmarshal DynamoDB record")
	}

	fmt.Println(dynamoRecord)
	fmt.Println(dynamoRecord.Change)
	fmt.Println(dynamoRecord.Change.NewImage)
	if dynamoRecord.Change.NewImage != nil {
		key, ok := dynamoRecord.Change.NewImage[filenameKey]
		if !ok {
			return errors.New("attribute 'filenameKey' not found in new image")
		}

		fileBody, err := s3Bucket.Get(ctx, key.String(), serviceProcessingBucket)
		if err != nil {
			return errors.Wrapf(err, "failed to get file %s from bucket %s", key.String(), serviceProcessingBucket)
		}
		defer fileBody.Close()

		if err = s3Bucket.Put(ctx, key.String(), serviceDictionaryBucket, fileBody.(io.Reader), "application/octet-stream"); err != nil {
			return errors.Wrapf(err, "failed to upload file %s to bucket %s", key.String(), serviceDictionaryBucket)
		}
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
