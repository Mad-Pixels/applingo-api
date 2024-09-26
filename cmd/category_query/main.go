package main

import (
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	validator "github.com/go-playground/validator/v10"
)

// service vars.
// ...

// system vars.
var (
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func main() {
	aws_lambda.Start(
		lambda.NewLambda(
			lambda.Config{},
			map[string]lambda.HandleFunc{
				"query": handleGet,
			},
		).Handle,
	)
}
