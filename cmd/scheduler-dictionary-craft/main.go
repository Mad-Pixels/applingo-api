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

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
)

const (
	lambdaWatchdog       = 240 * time.Second
	backoffOpenAIRequest = 15 * time.Second
	retriesOpenAIRequest = 1
	backoffBucketCheck   = 300 * time.Millisecond
	retriesBucketCheck   = 4
	defaultMaxWorkers    = 2
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

	timeout = utils.GetTimeout(lambdaTimeout, lambdaWatchdog)
)

func prepareWithDefaults(i *int, defaultValue int) *int {
	val := defaultValue
	if i == nil {
		return &val
	}
	if *i <= 0 || *i > defaultValue {
		return &val
	}
	return i
}

func init() {
	debug.SetGCPercent(500)

	gptClient = chatgpt.MustClient(
		httpclient.New().
			WithTimeout(timeout).
			WithMaxRetries(retriesOpenAIRequest, backoffOpenAIRequest).
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
	var (
		request     forge.RequestDictionaryCraft
		dynamoItems []applingoprocessing.SchemaItem
	)
	if err := serializer.UnmarshalJSON(record, &request); err != nil {
		return fmt.Errorf("failed to unmarshal request record: %w", err)
	}
	request.DictionariesCount = prepareWithDefaults(request.DictionariesCount, maxCraftDictionaries)
	request.MaxConcurrent = prepareWithDefaults(request.MaxConcurrent, maxCraftConurrent)

	craftResult, craftErrs := forge.CraftMultiple(ctx, &request, serviceForgeBucket, gptClient, s3Bucket)
	if len(craftErrs) > 0 {
		for _, err := range craftErrs {
			log.Error().Err(err).Msg("craft task was failed")
		}
	}
	for _, dictionary := range craftResult {
		if dictionary == nil {
			log.Error().Msg("dictionary is nil")
			continue
		}

		content, err := serializer.MarshalJSON(dictionary.GetWordsContainer())
		if err != nil {
			log.Error().Any("dictionary", *dictionary).Err(err).Msg("wrong dictionary format")
			continue
		}

		if err = s3Bucket.Put(
			ctx,
			utils.RecordToFileID(utils.GenerateDictionaryID(dictionary.GetDictionaryName(), dictionary.GetDictionaryAuthor())),
			serviceProcessingBucket,
			bytes.NewReader(content),
			cloud.ContentTypeJSON,
		); err != nil {
			log.Error().Any("dictionary", *dictionary).Err(err).Msg("upload dictionary to bucket failed")
			continue
		}
		if err = s3Bucket.WaitOrError(
			ctx,
			utils.RecordToFileID(utils.GenerateDictionaryID(dictionary.GetDictionaryName(), dictionary.GetDictionaryAuthor())),
			serviceProcessingBucket,
			retriesBucketCheck,
			backoffBucketCheck,
		); err != nil {
			log.Error().Any("dictionary", *dictionary).Err(err).Msg("failed check data in processing bucket")
			continue
		}

		dynamoItem := applingoprocessing.SchemaItem{
			Id: utils.GenerateDictionaryID(dictionary.GetDictionaryName(), dictionary.GetDictionaryAuthor()),

			// language info.
			Languages:   utils.JoinValues(dictionary.GetLanguageFrom().Name, dictionary.GetLanguageTo().Name),
			Level:       dictionary.GetLanguageLevel().String(),
			Subcategory: dictionary.GetSubcategory(),

			// dictionary info.
			Words:    dictionary.GetWordsCount(),
			Overview: dictionary.GetDictionaryOverview(),
			Author:   dictionary.GetDictionaryAuthor(),
			Name:     dictionary.GetDictionaryName(),

			// craft info.
			PromptCraft: utils.JoinValues(dictionary.GetPrompt(), string(dictionary.GetModel())),
			Description: dictionary.GetDictionaryDescription(),
			Topic:       dictionary.GetDictionaryTopic(),

			// internal info.
			Upload:  applingoprocessing.BoolToInt(false),
			Created: int(time.Now().Unix()),
			Reason:  "waiting for check",
		}
		dynamoItems = append(dynamoItems, dynamoItem)
	}
	if len(dynamoItems) > 0 {
		dynamoItems, err := applingoprocessing.BatchPutItems(dynamoItems)
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
