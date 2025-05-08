// Package main implements the Lambda API for managing user profiles.
// It handles creating and updating profile records in DynamoDB,
// with appropriate permission checks and validation.
package main

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/validator"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	awsRegion = os.Getenv("AWS_REGION")
	validate  *validator.Validator
	dbDynamo  *cloud.Dynamo
)

func init() {
	debug.SetGCPercent(500)
	validate = validator.New()

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	dbDynamo = cloud.NewDynamo(cfg)
}

func main() {
	lambda.Start(
		api.NewLambda(
			api.Config{
				EnableRequestLogging: true,
			},
			map[string]api.HandleFunc{
				// create profile
				"POST:/v1/profile": handleProfilePost,

				// patch profile data
				"PATCH:/v1/profile": handleProfilePatch,
			},
		).Handle,
	)
}
