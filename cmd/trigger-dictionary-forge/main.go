package main

import (
	"context"
	"encoding/json"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/validator"
	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	defaultTimeout    = 90 * time.Second
	defaultMaxWorkers = 5
)

var (
	awsRegion   = os.Getenv("AWS_REGION")
	openaiToken = os.Getenv("OPENAI_KEY")

	validate  *validator.Validator
	dbDynamo  *cloud.Dynamo
	s3Bucket  *cloud.Bucket
	gptClient *chatgpt.Client
)

func init() {
	debug.SetGCPercent(500)

	gptClient = chatgpt.MustClient(httpclient.New().WithTimeout(defaultTimeout), openaiToken)
	validate = validator.New()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
	dbDynamo = cloud.NewDynamo(cfg)
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
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
