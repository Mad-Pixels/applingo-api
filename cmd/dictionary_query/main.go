package main

import (
	"context"
	"os"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	awsRegion = os.Getenv("AWS_REGION")
	dbDynamo  *cloud.Dynamo
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	dbDynamo = cloud.NewDynamo(cfg)
}

func main() {
	aws_lambda.Start(
		lambda.NewLambda(
			lambda.Config{},
			map[string]lambda.HandleFunc{
				"query": handleDataQuery,
			},
		).Handle,
	)
}
