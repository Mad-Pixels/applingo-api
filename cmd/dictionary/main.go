package main

import (
	"context"
	"encoding/json"
	"github.com/Mad-Pixels/lingocards-api/pkg/amz"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
)

var (
	//
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	//serviceDictionaryDynamo = os.Getenv("SERVICE_DICTIONARY_DYNAMO")

	//
	awsRegion = os.Getenv("AWS_REGION")
	sess      *session.Session
)

type Response struct {
	Message string `json:"message"`
}

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(awsRegion)},
	}))
}

func putRequestHandler(ctx context.Context, data json.RawMessage) (any, error) {
	obj := amz.NewS3(sess)
	result, err := obj.PutRequest("fname", serviceDictionaryBucket, "json")
	if err != nil {
		return nil, err
	}

	return Response{
		Message: result,
	}, nil
}

func main() {
	aws_lambda.Start(
		lambda.NewLambda(
			map[string]lambda.HandleFunc{
				"putRequest": putRequestHandler,
			},
		).Handle,
	)
}
