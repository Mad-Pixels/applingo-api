package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Dynamo struct {
	client *dynamodb.Client
}

// NewDynamo creates a dynamo object.
func NewDynamo(cfg aws.Config) *Dynamo {
	return &Dynamo{
		client: dynamodb.NewFromConfig(cfg),
	}
}

// Put item to DynamoDB table.
func (d *Dynamo) Put(ctx context.Context, table string, item map[string]types.AttributeValue) error {
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &table,
		Item:      item,
	})
	return err
}
