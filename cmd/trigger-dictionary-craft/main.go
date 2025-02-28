package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt/dictionary_craft"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	dbDynamo  *cloud.Dynamo
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
	dbDynamo = cloud.NewDynamo(cfg)
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
		return errors.Wrap(err, "failed to upload dictionary to S3")
	}

	dynamoItem, err := applingoprocessing.PutItem(applingoprocessing.SchemaItem{
		Id:          utils.GenerateDictionaryID(craft.Meta.DictionaryName, craft.Info.Author),
		Languages:   craft.Meta.LanguageFrom + "-" + craft.Meta.LanguageTo,
		Description: craft.Meta.DictionaryDescription,
		Topic:       craft.Meta.DictionaryTopic,
		File:        craft.Meta.DictionaryName,
		Level:       craft.Meta.LanguageLevel,
		Overview:    craft.Info.Description,
		Prompt:      craft.Meta.Prompt,
		Name:        craft.Info.Name,
	})
	if err != nil {
		return errors.Wrap(err, "failed convert dictionary to DynamoDB item")
	}
	if err = dbDynamo.Put(
		ctx,
		applingoprocessing.TableSchema.TableName,
		dynamoItem,
		expression.AttributeNotExists(expression.Name(applingoprocessing.ColumnId)),
	); err != nil {
		return errors.Wrap(err, "failed to put DynamoDB item")
	}

	log.Info().Any("matadata", craft.Meta).Msg("dictionary was created successfully")
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
