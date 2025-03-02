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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/pkg/errors"
)

const (
	defaultBackoff    = 15 * time.Second
	defaultRetries    = 2
	defaultMaxWorkers = 5
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
		return errors.Wrap(err, "failed to unmarshal request record")
	}
	result, err := forge.CraftMultiple(
		ctx,
		&request,
		4,
		serviceForgeBucket,
		gptClient,
		s3Bucket,
		2,
	)
	if err != nil {
		return errors.Wrap(err[0], "failed to craft dictionary")
	}
	for _, item := range result {
		if item == nil {
			continue
		}

		content, err := serializer.MarshalJSON(item)
		if err != nil {
			continue
		}
		err = s3Bucket.Put(
			ctx,
			*item.Request.DictionaryFile,
			serviceProcessingBucket,
			bytes.NewReader(content),
			cloud.ContentTypeJSON,
		)
		if err != nil {
			continue
		}

		dynamoItem, err := applingoprocessing.PutItem(applingoprocessing.SchemaItem{
			Id:          utils.GenerateDictionaryID(item.Meta.Name, item.Meta.Author),
			Languages:   aws.ToString(item.Request.LanguageFrom) + "-" + aws.ToString(item.Request.LanguageTo),
			Description: aws.ToString(item.Request.DictionaryDescription),
			Topic:       aws.ToString(item.Request.DictionaryTopic),
			File:        aws.ToString(item.Request.DictionaryFile),
			Level:       aws.ToString(item.Request.LanguageLevel),
			Overview:    item.Meta.Description,
			Prompt:      aws.ToString(item.Request.Prompt),
			Name:        item.Meta.Name,
		})
		if err != nil {
			continue
			//return errors.Wrap(err, "failed convert dictionary to DynamoDB item")
		}
		if err = dbDynamo.Put(
			ctx,
			applingoprocessing.TableSchema.TableName,
			dynamoItem,
			expression.AttributeNotExists(expression.Name(applingoprocessing.ColumnId)),
		); err != nil {
			return errors.Wrap(err, "failed to put DynamoDB item")
		}
		log.Info().Any("matadata", item.Request).Msg("dictionary was created successfully")
	}

	// content, err := craft.Marshal()
	// if err != nil {
	// 	return errors.Wrap(err, "failed to marshal dictionary")
	// }
	// err = s3Bucket.Put(
	// 	ctx,
	// 	request.DictionaryName,
	// 	serviceProcessingBucket,
	// 	bytes.NewReader(content),
	// 	cloud.ContentTypeJSON,
	// )
	// if err != nil {
	// 	return errors.Wrap(err, "failed to upload dictionary to S3")
	// }

	// dynamoItem, err := applingoprocessing.PutItem(applingoprocessing.SchemaItem{
	// 	Id:          utils.GenerateDictionaryID(craft.Meta.DictionaryName, craft.Info.Author),
	// 	Languages:   craft.Meta.LanguageFrom + "-" + craft.Meta.LanguageTo,
	// 	Description: craft.Meta.DictionaryDescription,
	// 	Topic:       craft.Meta.DictionaryTopic,
	// 	File:        craft.Meta.DictionaryName,
	// 	Level:       craft.Meta.LanguageLevel,
	// 	Overview:    craft.Info.Description,
	// 	Prompt:      craft.Meta.Prompt,
	// 	Name:        craft.Info.Name,
	// })
	// if err != nil {
	// 	return errors.Wrap(err, "failed convert dictionary to DynamoDB item")
	// }
	// if err = dbDynamo.Put(
	// 	ctx,
	// 	applingoprocessing.TableSchema.TableName,
	// 	dynamoItem,
	// 	expression.AttributeNotExists(expression.Name(applingoprocessing.ColumnId)),
	// ); err != nil {
	// 	return errors.Wrap(err, "failed to put DynamoDB item")
	//}

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
