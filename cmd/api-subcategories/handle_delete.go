package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingosubcategory"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func handleDelete(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	id := generateSubcategoryID(string(*baseParams.GetStringPtr("code")), string(*baseParams.GetStringPtr("side")))
	result, err := dbDynamo.Get(ctx, applingosubcategory.TableName, map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	})
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: errors.Wrap(err, "failed to get item for deletion")}
	}
	if result.Item == nil {
		return nil, &api.HandleError{Status: http.StatusNotFound, Err: errors.New("item not found")}
	}

	if err := dbDynamo.Delete(ctx, applingosubcategory.TableName, map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: errors.Wrap(err, "failed to delete item")}
	}
	return openapi.DataResponseSuccess, nil
}
