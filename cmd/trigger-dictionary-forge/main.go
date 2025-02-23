package main

import (
	"bytes"
	"context"
	"encoding/csv"
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

// GPTDictionaryResponse структура для парсинга ответа от GPT
type GPTDictionaryResponse struct {
	Words []struct {
		Word        string `json:"word"`
		Translation string `json:"translation"`
		Hint        string `json:"hint"`
		Description string `json:"description"`
	} `json:"words"`
}

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
	log.Info().Any("data", request).Msg("dictionary craft parameters")

	prompt, err := preparePrompt(ctx, s3Bucket, request)
	if err != nil {
		return err
	}

	gptReq := &chatgpt.Request{
		Model: request.Model,
		Messages: []chatgpt.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	resp, err := gptClient.SendMessage(ctx, gptReq)
	if err != nil {
		return errors.Wrap(err, "failed to process ChatGPT request")
	}

	// Парсим JSON ответ от GPT
	var dictResponse GPTDictionaryResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &dictResponse); err != nil {
		return errors.Wrap(err, "failed to parse dictionary response")
	}

	// Создаем CSV из полученных данных
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Записываем данные в CSV
	for _, word := range dictResponse.Words {
		record := []string{
			word.Word,
			word.Translation,
			word.Hint,
			word.Description,
		}
		if err := writer.Write(record); err != nil {
			return errors.Wrap(err, "failed to write CSV record")
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return errors.Wrap(err, "failed to flush CSV writer")
	}

	// Сохраняем CSV в S3
	err = s3Bucket.Put(
		ctx,
		request.DictionaryName,
		os.Getenv("SERVICE_PROCESSING_BUCKET"),
		bytes.NewReader(buf.Bytes()),
		cloud.ContentTypeCSV,
	)
	if err != nil {
		return errors.Wrap(err, "failed to upload CSV to S3")
	}

	log.Info().
		Str("filename", request.DictionaryName).
		Msg("dictionary created successfully")

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
