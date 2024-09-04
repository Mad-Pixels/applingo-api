package main

import (
	"context"
	"github.com/Mad-Pixels/lingocards-api/pkg/amz"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event Event) (Response, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))},
	}))

	obj := amz.NewS3(sess)
	result, err := obj.PutRequest("fname", "lingocards.dictionary", "json")
	if err != nil {
		return Response{
			Message: "err.Error()",
		}, err
	}
	return Response{
		Message: result,
	}, nil
}

func main() {
	lambda.Start(handler)
}
