package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

func handler(ctx context.Context, event Event) (Response, error) {
	return Response{
		Message: "message",
	}, nil
}

func main() {
	lambda.Start(handler)
}
