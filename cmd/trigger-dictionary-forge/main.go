package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt/dictionary_craft"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
)

const (
	defaultTimeout    = 240 * time.Second
	defaultBackoff    = 15 * time.Second
	defaultRetries    = 2
	defaultMaxWorkers = 5
)

var (
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	serviceForgeBucket      = os.Getenv("SERVICE_FORGE_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")
	openaiToken             = os.Getenv("OPENAI_KEY")

	gptClient *chatgpt.Client
	s3Bucket  *cloud.Bucket
)

func init() {
	debug.SetGCPercent(500)

	gptClient = chatgpt.MustClient(
		httpclient.New().
			WithTimeout(defaultTimeout).
			WithMaxRetries(defaultRetries, defaultBackoff).
			WithRetryCondition(func(statusCode int, responseBody string) bool {
				return statusCode >= 500 && statusCode < 600
			}),
		openaiToken,
	)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	var request dictionary_craft.Request
	if err := request.Unmarshal(record); err != nil {
		return errors.Wrap(err, "failed to unmarshal request record")
	}
	craft, err := dictionary_craft.Craft(ctx, &request, serviceForgeBucket, gptClient, s3Bucket)
	if err != nil {
		return errors.Wrap(err, "failed to craft dictionary")
	}
	log.Info().Any("matadata", craft.Meta).Msg("dictionary was created successfully")

	content, err := craft.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal dictionary")
	}
	err = s3Bucket.Put(
		ctx,
		request.DictionaryName,
		serviceProcessingBucket,
		bytes.NewReader(content),
		cloud.ContentTypeJSON,
	)
	if err != nil {
		return errors.Wrap(err, "failed to upload CSV to S3")
	}
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
