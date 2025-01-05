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
	openaiPrompt            = os.Getenv("OPENAI_PROMPT")
	openaiModel             = os.Getenv("OPENAI_MODEL")
	openaiKey               = os.Getenv("OPENAI_KEY")
	awsRegion               = os.Getenv("AWS_REGION")

	validate *validator.Validator
	s3Bucket *cloud.Bucket
)

func init() {
	debug.SetGCPercent(500)
	validate = validator.New()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	prompt, err := s3Bucket.Read(ctx, openaiPrompt, serviceForgeBucket)
	if err != nil {
		return err
	}

	request := GPTRequest{
		Model: openaiModel,
		Messages: []Message{
			{
				Role:    "user",
				Content: string(prompt),
			},
		},
		Temperature: temperature,
	}
	payload, err := serializer.MarshalJSON(request)
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
