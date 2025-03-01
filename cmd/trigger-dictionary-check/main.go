package main

import (
	"context"
	"encoding/json"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	awsRegion = os.Getenv("AWS_REGION")
)

func init() {
	debug.SetGCPercent(500)

	_, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
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

	log.Error().Any("item", item).Msg("item")
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
