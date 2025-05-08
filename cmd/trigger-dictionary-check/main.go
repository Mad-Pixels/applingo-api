// Package main implements a Lambda function to handle DynamoDB REMOVE events
// and delete associated files from the S3 dictionary bucket.
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
	defaultMaxWorkers = 20
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")

	s3Bucket *cloud.Bucket
)

func init() {
	debug.SetGCPercent(500)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
}

// handler processes DynamoDB REMOVE events and deletes corresponding files from S3.
func handler(ctx context.Context, _ zerolog.Logger, record json.RawMessage) error {
	var dynamoDBEvent events.DynamoDBEventRecord
	if err := serializer.UnmarshalJSON(record, &dynamoDBEvent); err != nil {
		return fmt.Errorf("failed to unmarshal request record: %w", err)
	}

	if dynamoDBEvent.EventName == "REMOVE" {
		if err := remove(ctx, dynamoDBEvent); err != nil {
			return fmt.Errorf("failed to delete file from bucket: %w", err)
		}
	}

	return nil
}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{
				MaxWorkers: defaultMaxWorkers,
			},
			handler,
		).Handle,
	)
}
