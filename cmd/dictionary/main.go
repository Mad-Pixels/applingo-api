package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	lambda2 "github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/pkg/amz"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceDictionaryDynamo = os.Getenv("SERVICE_DICTIONARY_DYNAMO")
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
}

func putRequestHandler(ctx context.Context, data json.RawMessage) (any, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("invalid request format: %w", err)
	}

	l := lambda2.MustLambda(nil)
	obj := amz.NewS3(l.AwsSes)
	result, err := obj.PutRequest("fname", "lingocards.dictionary", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to put request: %w", err)
	}

	return Response{
		Message: result,
	}, nil
}

func main() {
	handlers := map[string]lambda2.HandleFunc{
		"putRequest": putRequestHandler,
	}
	l := lambda2.MustLambda(handlers)
	lambda.Start(l.Route)
}
