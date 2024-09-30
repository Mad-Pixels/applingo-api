package main

import (
	"github.com/Mad-Pixels/lingocards-api/pkg/api"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"runtime/debug"
)

var (
	validate *validator.Validate
)

func init() {
	debug.SetGCPercent(500)
	validate = validator.New()
}

func main() {
	lambda.Start(
		api.NewLambda(
			api.Config{},
			map[string]api.HandleFunc{
				"query": handleGet,
			},
		).Handle,
	)
}
