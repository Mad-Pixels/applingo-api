package main

import (
	"runtime/debug"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	debug.SetGCPercent(200)
}

func main() {
	lambda.Start(
		api.NewLambda(
			api.Config{
				EnableRequestLogging: true,
			},
			map[string]api.HandleFunc{
				// list
				"GET:/v1/schema": handleSchemaGet,
			},
		).Handle,
	)
}
