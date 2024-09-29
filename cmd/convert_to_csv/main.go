package main

import (
	"context"
	"encoding/csv"
	"fmt"
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

	awsRegion = os.Getenv("AWS_REGION")
	s3Bucket  *cloud.Bucket
	dbDynamo  *cloud.Dynamo
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
			return fmt.Errorf("error getting object %s/%s: %v", serviceProcessingBucket, key, err)
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
				return fmt.Errorf("error reading CSV: %v", err)
			}
			if err := csvWriter.Write(record); err != nil {
				return fmt.Errorf("error writing CSV: %v", err)
			}
		}
		csvWriter.Flush()

		newKey := strings.TrimSuffix(key, ".csv") + "_converted.csv"
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
