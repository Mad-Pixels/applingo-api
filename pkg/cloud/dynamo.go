package cloud

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dynamo struct {
	client *dynamodb.Client
}

// NewDynamo creates a dynamo object.
func NewDynamo(cfg aws.Config) *dynamo {
	return &dynamo{
		client: dynamodb.NewFromConfig(cfg),
	}
}
