package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")
	logger                  = lambda.InitLogger()

	s3Bucket *cloud.Bucket
	dbDynamo *cloud.Dynamo
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
	dbDynamo = cloud.NewDynamo(cfg)

	debug.SetGCPercent(500)
}

func handler(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		key := record.S3.Object.Key

		reader, err := s3Bucket.Get(ctx, key, serviceProcessingBucket)
		if err != nil {
			logger.Error().Err(err).Str("bucket", serviceProcessingBucket).Str("key", key).Msg("cannot get object from bucket")
			return errors.New("exit with error")
		}
		defer reader.Close()

		csvData := &strings.Builder{}
		csvWriter := csv.NewWriter(csvData)
		csvReader := csv.NewReader(reader)
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Error().Err(err).Str("bucket", serviceProcessingBucket).Str("key", key).Msg("error reading CSV")
				return errors.New("exit with error")
			}
			if err := csvWriter.Write(record); err != nil {
				logger.Error().Err(err).Str("bucket", serviceProcessingBucket).Str("key", key).Msg("error writing CSV")
				return errors.New("exit with error")
			}
		}
		csvWriter.Flush()

		newKey := uuid.New().String() + ".csv"
		err = s3Bucket.Put(ctx, newKey, serviceDictionaryBucket, strings.NewReader(csvData.String()), "text/csv")
		if err != nil {
			return fmt.Errorf("unable to upload object to %s/%s: %v", serviceDictionaryBucket, newKey, err)
		}
	}
	return nil
}

func main() {
	aws_lambda.Start(handler)
}
