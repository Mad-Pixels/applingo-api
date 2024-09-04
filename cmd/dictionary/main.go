package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mad-Pixels/lingocards-api/internal"
	"github.com/Mad-Pixels/lingocards-api/pkg/amz"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
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

	l := internal.MustLambda(nil)
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
	handlers := map[string]internal.HandleFunc{
		"putRequest": putRequestHandler,
	}
	l := internal.MustLambda(handlers)
	lambda.Start(l.Route)
}
