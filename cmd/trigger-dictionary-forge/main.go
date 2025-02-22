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
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
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
	openaiKey               = os.Getenv("OPENAI_KEY")
	awsRegion               = os.Getenv("AWS_REGION")

	httpCli  *httpclient.ClientWrapper
	validate *validator.Validator
	s3Bucket *cloud.Bucket
)

func init() {
	debug.SetGCPercent(500)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	httpCli = httpclient.New().WithTimeout(time.Duration(requestTimeout) * time.Second)
	s3Bucket = cloud.NewBucket(cfg)
	validate = validator.New()
}

func promptPrepare(ctx context.Context, req ForgeRequest) (string, error) {
	tpl, err := s3Bucket.Read(ctx, req.OpenAIPromptName, serviceForgeBucket)
	if err != nil {
		return "", errors.Wrap(err, "failed get prompt")
	}
	content, err := utils.Template(string(tpl), req)
	if err != nil {
		return "", errors.Wrap(err, "cannot prepare promt content")
	}
	return content, nil
}

func sendGPTResponse(ctx context.Context, model string, content string) (string, error) {
	req := GPTRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: content,
			},
		},
		Temperature: temperature,
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + openaiKey,
	}
	payload, err := serializer.MarshalJSON(req)
	if err != nil {
		return "", errors.Wrap(err, "openai request serialization failed")
	}

	body, err := httpCli.Post(ctx, url, string(payload), headers)
	if err != nil {
		return "", errors.Wrap(err, "failed to get response from openai")
	}
	var resp GPTResponse
	if err := serializer.UnmarshalJSON([]byte(body), &resp); err != nil {
		return "", errors.Wrap(err, "failed to parse response from openai")
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("got empty response")
	}
	return resp.Choices[0].Message.Content, nil
}

func handler(ctx context.Context, log zerolog.Logger, body json.RawMessage) error {
	var req ForgeRequest
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return errors.Wrap(err, "bad request, cannot deserialize request data")
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return errors.Wrap(err, "bad request, validateion failed")
	}

	prompt, err := promptPrepare(ctx, req)
	if err != nil {
		return err
	}
	content, err := sendGPTResponse(ctx, req.OpenAIModelName, prompt)
	if err != nil {
		return errors.Wrap(err, "failed process openapi request")
	}
	table, err := utils.CSV(content)
	if err != nil {
		return errors.Wrap(err, "failed process CSV from response body")
	}

	err = s3Bucket.Put(ctx, req.DictionaryName, serviceProcessingBucket, bytes.NewReader(table), cloud.ContentTypeCSV)
	if err != nil {
		return errors.Wrap(err, "failed to upload CSV to S3")
	}
	log.Info().
		Str("filename", req.DictionaryName).
		Msg("dictionary was created")
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
