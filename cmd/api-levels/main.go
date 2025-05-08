// Package main implements the Lambda handler for retrieving available language learning levels.
// It exposes the levels endpoint and returns a predefined list of CEFR levels.
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
				"GET:/v1/levels": handleLevelsGet,
			},
		).Handle,
	)
}
