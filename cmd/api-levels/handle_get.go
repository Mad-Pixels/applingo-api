package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingolevel"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/rs/zerolog"
)

const pageLimit = 6

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	scanInput := dbDynamo.BuildScanInput(applingolevel.TableName, pageLimit, nil)
	result, err := dbDynamo.Scan(ctx, applingolevel.TableName, scanInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	var items []applingoapi.LevelItemV1
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &items); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	return openapi.DataResponseLevels(applingoapi.LevelsData{
		Items: items,
	}), nil
}
