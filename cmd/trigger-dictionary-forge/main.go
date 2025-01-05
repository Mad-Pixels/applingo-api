package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/validator"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	url = "https://api.openai.com/v1/chat/completions"

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
		return errors.Wrap(err, "failed get prompt")
	}
	httpcli := httpclient.New().WithTimeout(time.Duration(requestTimeout) * time.Second)

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
		return errors.Wrap(err, "request serialization error")
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + openaiKey,
	}

	responseRaw, err := httpcli.Post(ctx, url, string(payload), headers)
	if err != nil {
		return errors.Wrap(err, "failed get response from OpenAI")
	}
	var response GPTResponse
	if err := serializer.UnmarshalJSON([]byte(responseRaw), &response); err != nil {
		return errors.Wrap(err, "failed to parse GPT response")
	}
	if len(response.Choices) == 0 {
		return errors.New("empty response form GPT")
	}

	cvsData, err := toCSV(response.Choices[0].Message.Content)
	if err != nil {
		return errors.Wrap(err, "failed to convert GPT response to CSV")
	}
	err = s3Bucket.Put(ctx, "filename", serviceProcessingBucket, bytes.NewReader(cvsData), cloud.ContentTypeCSV)
	if err != nil {
		return errors.Wrap(err, "failed to upload CSV to S3")
	}

	log.Info().
		Str("filename", "set filename").
		Msg("successfully processed and uploaded GPT response")
	return nil
}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{MaxWorkers: 5},
			handler,
		).Handle,
	)
}
