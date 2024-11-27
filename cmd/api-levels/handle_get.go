package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingolevel"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

const pageLimit = 6

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	var items []map[string]types.AttributeValue
	scanInput := dbDynamo.BuildScanInput(applingolevel.TableName, pageLimit, nil)
	result, err := dbDynamo.Scan(ctx, applingolevel.TableName, scanInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	items = result.Items

	var (
		wg      sync.WaitGroup
		itemsCh = make(chan applingoapi.LevelItemV1, len(items))
	)
	response := applingoapi.LevelsData{}

	for _, item := range items {
		wg.Add(1)
		go func(item map[string]types.AttributeValue) {
			defer wg.Done()

			var level applingoapi.LevelItemV1
			if err := attributevalue.UnmarshalMap(item, &level); err != nil {
				logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
				return
			}
			itemsCh <- level
		}(item)
	}
	go func() {
		wg.Wait()
		close(itemsCh)
	}()

	for item := range itemsCh {
		response.Items = append(response.Items, item)
	}
	return openapi.DataResponseLevels(response), nil
}
