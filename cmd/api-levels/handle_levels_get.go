package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingolevel"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/rs/zerolog"
)

const pageLimit = 6

func handleLevelsGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	scanInput := dbDynamo.BuildScanInput(applingolevel.TableName, pageLimit, nil)
	result, err := dbDynamo.Scan(ctx, applingolevel.TableName, scanInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	var dynamoItems []applingolevel.SchemaItem
	if err := attributevalue.UnmarshalListOfMaps(result.Items, &dynamoItems); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	var items []applingoapi.LevelItemV1
	for _, item := range dynamoItems {
		items = append(items, applingoapi.LevelItemV1{
			Code:  item.Code,
			Level: item.Level,
		})
	}
	return openapi.DataResponseLevels(applingoapi.LevelsData{
		Items: items,
	}), nil
}
