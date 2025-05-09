// Package main implements the API Lambda for reporting runtime errors from the device.
// It provides the report endpoint for sending structured error reports to an S3 bucket.
package main

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/validator"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	serviceErrorsBucket = os.Getenv("SERVICE_ERRORS_BUCKET")
	awsRegion           = os.Getenv("AWS_REGION")

	validate *validator.Validator
	s3Bucket *cloud.Bucket
)

func init() {
	debug.SetGCPercent(500)
	validate = validator.New()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
}

func main() {
	lambda.Start(
		api.NewLambda(
			api.Config{
				EnableRequestLogging: true,
			},
			map[string]api.HandleFunc{
				// item
				"POST:/v1/report": handleReportPost,
			},
		).Handle,
	)
}
