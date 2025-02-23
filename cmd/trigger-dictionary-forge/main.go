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
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/Mad-Pixels/applingo-api/pkg/validator"
	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
)

const (
	defaultTimeout    = 90 * time.Second
	defaultMaxWorkers = 5
)

var (
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	serviceForgeBucket      = os.Getenv("SERVICE_FORGE_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")
	openaiToken             = os.Getenv("OPENAI_KEY")

	validate  *validator.Validator
	gptClient *chatgpt.Client
	dbDynamo  *cloud.Dynamo
	s3Bucket  *cloud.Bucket
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

func preparePrompt(ctx context.Context, bucket *cloud.Bucket, req Request) (string, error) {
	tpl, err := bucket.Read(ctx, req.Prompt, serviceForgeBucket)
	if err != nil {
		return "", errors.Wrap(err, "failed to get prompt template")
	}

	content, err := utils.Template(string(tpl), req)
	if err != nil {
		return "", errors.Wrap(err, "failed to prepare prompt content")
	}

	return content, nil
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	var request Request
	if err := serializer.UnmarshalJSON(record, &request); err != nil {
		return errors.Wrap(err, "failed to unmarshal request record")
	}

	if err := request.Update(ctx, s3Bucket, serviceForgeBucket); err != nil {
		return errors.Wrap(err, "failed to update request data")
	}

	prompt, err := preparePrompt(ctx, s3Bucket, request)
	if err != nil {
		return err
	}
	log.Info().Any("data", prompt).Msg("result")
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
