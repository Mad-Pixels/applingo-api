package main

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/go-playground/validator/v10"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")

	awsRegion = os.Getenv("AWS_REGION")
	validate  *validator.Validate
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
	validate = validator.New()

	debug.SetGCPercent(500)
}

func main() {
	aws_lambda.Start(
		lambda.NewLambda(
			lambda.Config{},
			map[string]lambda.HandleFunc{
				"query":        handleDataQuery,
				"download_url": handleDownloadUrl,
			},
		).Handle,
	)
}
