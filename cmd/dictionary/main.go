package main

import (
	"os"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator/v10"
)

var (
	// service vars.
	//serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	//serviceDictionaryDynamo = os.Getenv("SERVICE_DICTIONARY_DYNAMO")

	// system vars.
	awsRegion = os.Getenv("AWS_REGION")
	validate  *validator.Validate
	sess      *session.Session
)

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(awsRegion)},
	}))
	validate = validator.New()
}

func main() {
	aws_lambda.Start(
		lambda.NewLambda(
			map[string]lambda.HandleFunc{
				"presign": handlePresign,
			},
		).Handle,
	)
}
