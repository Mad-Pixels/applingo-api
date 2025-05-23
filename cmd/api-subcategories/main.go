// Package main implements the AWS Lambda handler for the subcategories endpoint.
// It sets up routing and logging configuration using the API framework.
package main

import (
	"runtime/debug"

	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	debug.SetGCPercent(500)
}

func main() {
	lambda.Start(
		api.NewLambda(
			api.Config{
				EnableRequestLogging: true,
			},
			map[string]api.HandleFunc{
				// list
				"GET:/v1/subcategories": handleSubcategoriesGet,
			},
		).Handle,
	)
}
