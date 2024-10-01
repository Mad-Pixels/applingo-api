package main

import (
	"context"
	"encoding/json"
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

var (
	servicePutScvQueueUrl = os.Getenv("SERVICE_PUT_CSV_QUEUE_URL")
	awsRegion             = os.Getenv("AWS_REGION")

	sqsQueue *cloud.Queue
)

func init() {
	debug.SetGCPercent(500)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	sqsQueue = cloud.NewQueue(cfg)
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	var dynamoRecord events.DynamoDBEventRecord
	if err := serializer.UnmarshalJSON(record, &dynamoRecord); err != nil {
		return errors.Wrap(err, "failed to unmarshal DynamoDB record")
	}
	payload, err := serializer.MarshalJSON(dynamoRecord)
	if err != nil {
		return errors.Wrap(err, "failed to marshal DynamoDB record")
	}

	_, err = sqsQueue.SendMessage(ctx, servicePutScvQueueUrl, string(payload))
	if err != nil {
		return errors.Wrap(err, "failed to send message to SQS")
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
