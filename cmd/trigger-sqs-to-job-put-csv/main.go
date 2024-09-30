package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/go-playground/validator/v10"
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")

	validate *validator.Validate
	s3Bucket *cloud.Bucket
	dbDynamo *cloud.Dynamo
)

func init() {
	debug.SetGCPercent(500)

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
	dbDynamo = cloud.NewDynamo(cfg)
	validate = validator.New()
}

const filenameKey = "filename"

func Handler(ctx context.Context, event events.DynamoDBEvent) error {
	fmt.Println("serviceDictionaryBucket", serviceDictionaryBucket)
	fmt.Println("serviceProcessingBucket", serviceProcessingBucket)
	//for _, record := range event.Records {
	//	if record.Change.NewImage != nil {
	//		key, ok := record.Change.NewImage[filenameKey]
	//		if !ok {
	//			fmt.Println("Attribute 'Name' not found in new image")
	//			return errors.New("")
	//		}
	//
	//		fileBody, err := s3Bucket.Get(ctx, key.String(), serviceProcessingBucket)
	//		if err != nil {
	//			log.Printf("Failed to get file from bucket %s with key %s: %v", serviceProcessingBucket, key.String(), err)
	//			return err
	//		}
	//		defer fileBody.Close()
	//
	//		err = s3Bucket.Put(ctx, key.String(), serviceDictionaryBucket, fileBody, "application/octet-stream")
	//		if err != nil {
	//			log.Printf("Failed to upload file to bucket %s with key %s: %v", serviceDictionaryBucket, key.String(), err)
	//			return err
	//		}
	//
	//		fmt.Printf("Successfully moved file %s from %s to %s\n", key.String(), serviceProcessingBucket, serviceDictionaryBucket)
	//	}
	//}
	return errors.New("message with error")
	//return nil
}

func main() {
	lambda.Start(Handler)
}
