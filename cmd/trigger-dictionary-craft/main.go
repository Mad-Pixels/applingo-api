package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/forge"
	"github.com/Mad-Pixels/applingo-api/pkg/httpclient"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/rs/zerolog"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	defaultBackoff       = 15 * time.Second
	defaultRetries       = 2
	defaultMaxWorkers    = 5
	maxCraftConurrent    = 4
	maxCraftDictionaries = 4
)

var (
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	serviceForgeBucket      = os.Getenv("SERVICE_FORGE_BUCKET")
	lambdaTimeout           = os.Getenv("LAMBDA_TIMEOUT_SECONDS")
	awsRegion               = os.Getenv("AWS_REGION")
	openaiToken             = os.Getenv("OPENAI_KEY")

	gptClient *chatgpt.Client
	dbDynamo  *cloud.Dynamo
	s3Bucket  *cloud.Bucket

	timeout = utils.GetTimeout(lambdaTimeout, 240*time.Second)
)

func init() {
	debug.SetGCPercent(500)

	gptClient = chatgpt.MustClient(
		httpclient.New().
			WithTimeout(timeout).
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
	var request forge.RequestDictionaryCraft
	if err := serializer.UnmarshalJSON(record, &request); err != nil {
		return fmt.Errorf("failed to unmarshal request record: %w", err)
	}
	if request.MaxConcurrent == 0 || request.MaxConcurrent > maxCraftConurrent {
		request.MaxConcurrent = maxCraftConurrent
	}
	if request.DictionariesCount < 0 {
		request.DictionariesCount = 1
	}
	if request.DictionariesCount > maxCraftDictionaries {
		request.DictionariesCount = maxCraftDictionaries
	}

	var (
		schemaItems []applingoprocessing.SchemaItem
	)
	result, craftErrs := forge.CraftMultiple(ctx, &request, serviceForgeBucket, gptClient, s3Bucket)
	if len(craftErrs) > 0 {
		for _, err := range craftErrs {
			log.Error().Err(err).Msg("craft task was failed")
		}
	}

	for _, dictionary := range result {
		if dictionary == nil {
			log.Error().Msg("dictionary is nil")
			continue
		}

		content, err := serializer.MarshalJSON(dictionary)
		if err != nil {
			log.Error().Any("dictionary", *dictionary).Err(err).Msg("wrong dictionary format")
		}
		if err = s3Bucket.Put(
			ctx,
			dictionary.Request.GetDictionaryFile(),
			serviceProcessingBucket,
			bytes.NewReader(content),
			cloud.ContentTypeJSON,
		); err != nil {
			log.Error().Any("dictionary", *dictionary).Err(err).Msg("upload dictionary to bucker failed")
			continue
		}

		schemaItem := applingoprocessing.SchemaItem{
			Id:          utils.GenerateDictionaryID(dictionary.Meta.Name, dictionary.Meta.Author),
			Languages:   aws.ToString(dictionary.Request.LanguageFrom) + "-" + aws.ToString(dictionary.Request.LanguageTo),
			Description: aws.ToString(dictionary.Request.DictionaryDescription),
			Topic:       aws.ToString(dictionary.Request.DictionaryTopic),
			Level:       aws.ToString(dictionary.Request.LanguageLevel),
			PromptCraft: aws.ToString(dictionary.Request.Prompt),
			File:        dictionary.Request.GetDictionaryFile(),
			Upload:      applingoprocessing.BoolToInt(false),
			Subcategory: dictionary.Request.Subcategory(),
			Words:       dictionary.Request.WordsCount,
			Overview:    dictionary.Meta.Description,
			Author:      dictionary.Meta.Author,
			Name:        dictionary.Meta.Name,
			Created:     int(time.Now().Unix()),
			Reason:      "wait for check",
			PromptCheck: "",
			Score:       0,
		}
		schemaItems = append(schemaItems, schemaItem)
	}

	if len(schemaItems) > 0 {
		dynamoItems, err := applingoprocessing.BatchPutItems(schemaItems)
		if err != nil {
			log.Error().Err(err).Msg("failed to prepare batch items")
			return fmt.Errorf("failed to prepare batch items: %w", err)
		}
		if err = dbDynamo.BatchWrite(ctx, applingoprocessing.TableSchema.TableName, dynamoItems); err != nil {
			log.Error().Err(err).Msg("failed to batch write items to DynamoDB")
			return fmt.Errorf("failed to batch write items: %w", err)
		}
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
