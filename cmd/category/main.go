package main

import (
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

var (
	// service vars.
	// ...

	// system vars.
	validate *validator.Validate
)

func init() {
	validate = validator.New()
}

func main() {
	aws_lambda.Start(
		lambda.NewLambda(
			map[string]lambda.HandleFunc{
				"get": handleGet,
			},
		).Handle,
	)
}
