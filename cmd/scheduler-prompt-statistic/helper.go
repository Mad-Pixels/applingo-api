package main

import (
	"context"
	"fmt"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func fetchProcessingItems(ctx context.Context, dynamo *cloud.Dynamo, table string) ([]applingoprocessing.SchemaItem, error) {
	var allItems []applingoprocessing.SchemaItem
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		scanInput := dynamo.BuildScanInput(table, 100, lastEvaluatedKey)
		result, err := dynamo.Scan(ctx, table, scanInput)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		var items []applingoprocessing.SchemaItem
		err = attributevalue.UnmarshalListOfMaps(result.Items, &items)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %v", err)
		}

		allItems = append(allItems, items...)

		if result.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = result.LastEvaluatedKey
	}

	return allItems, nil
}
