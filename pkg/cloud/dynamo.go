package cloud

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

// Common errors
var (
	ErrDynamoEmptyTable = errors.New("empty table name")
	ErrDynamoEmptyKey   = errors.New("empty key")
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

// Dynamo represents a DynamoDB client for database operations.
type Dynamo struct {
	client *dynamodb.Client
}

// NewDynamo creates a new instance of DynamoDB client.
func NewDynamo(cfg aws.Config) *Dynamo {
	return &Dynamo{
		client: dynamodb.NewFromConfig(cfg),
	}
}

// validateTable checks if table name is not empty.
func validateTable(table string) error {
	if table == "" {
		return ErrDynamoEmptyTable
	}
	return nil
}

// validateKey checks if key is not empty.
func validateKey(key map[string]types.AttributeValue) error {
	if len(key) == 0 {
		return ErrDynamoEmptyKey
	}
	return nil
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
		return nil, errors.Wrap(err, "failed to build expression")
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

// BuildScanInput creates a dynamodb.ScanInput based on the provided fields and conditions.
func (d *Dynamo) BuildScanInput(table string, limit int32, exclusiveStartKey map[string]types.AttributeValue) *dynamodb.ScanInput {
	input := &dynamodb.ScanInput{
		TableName:         aws.String(table),
		Limit:             &limit,
		ExclusiveStartKey: exclusiveStartKey,
	}
	return input
}

// Put adds or updates an item in the DynamoDB table.
func (d *Dynamo) Put(ctx context.Context, table string, item map[string]types.AttributeValue, condition expression.ConditionBuilder) error {
	if err := validateTable(table); err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      item,
	}
	if condition.IsSet() {
		expr, err := expression.NewBuilder().WithCondition(condition).Build()
		if err != nil {
			return fmt.Errorf("failed to build condition expression: %w", err)
		}
		input.ConditionExpression = expr.Condition()
		input.ExpressionAttributeNames = expr.Names()
		input.ExpressionAttributeValues = expr.Values()
	}
	_, err := d.client.PutItem(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to put item")
	}
	return nil
}

// Get retrieves an item from DynamoDB table by its key.
func (d *Dynamo) Get(ctx context.Context, table string, key map[string]types.AttributeValue) (*dynamodb.GetItemOutput, error) {
	if err := validateTable(table); err != nil {
		return nil, err
	}
	if err := validateKey(key); err != nil {
		return nil, err
	}

	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       key,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get item")
	}
	return result, nil
}

// Query executes a query operation on DynamoDB table.
func (d *Dynamo) Query(ctx context.Context, table string, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	if err := validateTable(table); err != nil {
		return nil, err
	}

	input.TableName = aws.String(table)
	result, err := d.client.Query(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return result, nil
}

// Delete removes an item from DynamoDB table by key.
func (d *Dynamo) Delete(ctx context.Context, table string, key map[string]types.AttributeValue) error {
	if err := validateTable(table); err != nil {
		return err
	}
	if err := validateKey(key); err != nil {
		return err
	}

	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key:       key,
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete item")
	}
	return nil
}

// Update modifies an existing item in the DynamoDB table.
func (d *Dynamo) Update(ctx context.Context, table string, key map[string]types.AttributeValue, update expression.UpdateBuilder, condition expression.ConditionBuilder) error {
	if err := validateTable(table); err != nil {
		return err
	}
	if err := validateKey(key); err != nil {
		return err
	}

	builder := expression.NewBuilder().WithUpdate(update)
	if condition.IsSet() {
		builder = builder.WithCondition(condition)
	}
	expr, err := builder.Build()
	if err != nil {
		return errors.Wrap(err, "failed to build update expression")
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(table),
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}
	if expr.Condition() != nil {
		input.ConditionExpression = expr.Condition()
	}
	_, err = d.client.UpdateItem(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to update item")
	}
	return nil
}

// Scan executes a scan operation on DynamoDB table.
func (d *Dynamo) Scan(ctx context.Context, table string, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	if err := validateTable(table); err != nil {
		return nil, err
	}

	input.TableName = aws.String(table)
	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute scan")
	}
	return result, nil
}
