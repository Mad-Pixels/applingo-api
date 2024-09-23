package cloud

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// QueryInput represents the input for a DynamoDB query.
type QueryInput struct {
	IndexName         string
	KeyCondition      expression.KeyConditionBuilder
	FilterCondition   expression.ConditionBuilder
	ProjectionFields  []string
	Limit             int32
	ScanForward       bool
	ExclusiveStartKey map[string]types.AttributeValue
}

type Dynamo struct {
	client *dynamodb.Client
}

// NewDynamo creates a dynamo object.
func NewDynamo(cfg aws.Config) *Dynamo {
	return &Dynamo{
		client: dynamodb.NewFromConfig(cfg),
	}
}

// BuildQueryInput creates a dynamodb.QueryInput based on the provided QueryInput.
func (d *Dynamo) BuildQueryInput(input QueryInput) (*dynamodb.QueryInput, error) {
	builder := expression.NewBuilder().WithKeyCondition(input.KeyCondition)

	if input.FilterCondition.IsSet() {
		builder = builder.WithFilter(input.FilterCondition)
	}

	if len(input.ProjectionFields) > 0 {
		projBuilder := expression.ProjectionBuilder{}
		for _, field := range input.ProjectionFields {
			projBuilder = projBuilder.AddNames(expression.Name(field))
		}
		builder = builder.WithProjection(projBuilder)
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, err
	}

	queryInput := &dynamodb.QueryInput{
		IndexName:                 &input.IndexName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     &input.Limit,
		ScanIndexForward:          &input.ScanForward,
		ExclusiveStartKey:         input.ExclusiveStartKey,
	}

	if expr.Filter() != nil {
		queryInput.FilterExpression = expr.Filter()
	}

	if expr.Projection() != nil {
		queryInput.ProjectionExpression = expr.Projection()
	}

	return queryInput, nil
}

// Put item to DynamoDB table.
func (d *Dynamo) Put(ctx context.Context, table string, item map[string]types.AttributeValue) error {
	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &table,
		Item:      item,
	})
	return err
}

// Query executes a query operation on DynamoDB table.
func (d *Dynamo) Query(ctx context.Context, table string, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	input.TableName = &table
	return d.client.Query(ctx, input)
}

// Delete an item from DynamoDB table by key.
func (d *Dynamo) Delete(ctx context.Context, table string, key map[string]types.AttributeValue) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &table,
		Key:       key,
	})
	return err
}
