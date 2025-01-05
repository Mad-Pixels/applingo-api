package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/validator"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
)

var (
	requestTimeout = 90
	temperature    = 0.7
)

var (
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	serviceForgeBucket      = os.Getenv("SERVICE_FORGE_BUCKET")
	openapiModel            = os.Getenv("OPENAPI_MODEL")
	openapiKey              = os.Getenv("OPENAPI_KEY")
	awsRegion               = os.Getenv("AWS_REGION")

	validate   *validator.Validator
	s3Bucket   *cloud.Bucket
	httpClient *http.Client
)

func init() {
	debug.SetGCPercent(500)
	validate = validator.New()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
	httpClient = &http.Client{Timeout: time.Duration(requestTimeout) * time.Second}
}

// TODO: get content from S3 file as string.
func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	request := GPTRequest{
		Model: openapiModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: "",
			},
		},
		Temperature: temperature,
	}
	data, err := serializer.MarshalJSON(request)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{MaxWorkers: 5},
			handler,
		).Handle,
	)
}
