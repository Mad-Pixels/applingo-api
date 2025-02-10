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
				"GET:/v1/dictionaries":    handleGet,
				"POST:/v1/dictionaries":   handlePost,
				"DELETE:/v1/dictionaries": handleDelete,
				"PATCH:/v1/dictionaries":  handlePatchStatistic,
			},
		).Handle,
	)
}
