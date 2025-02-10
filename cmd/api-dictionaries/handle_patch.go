package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func handlePatchStatistic(ctx context.Context, _ zerolog.Logger, rawMessage json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	params := applingoapi.PatchStatisticDictionariesV1Params{
		Name:        baseParams.GetStringDefault("name", ""),
		Author:      baseParams.GetStringDefault("author", ""),
		Subcategory: baseParams.GetStringDefault("subcategory", ""),
	}
	if err := validate.ValidateStruct(&params); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	var req applingoapi.RequestPatchStatisticDictionariesV1
	if err := serializer.UnmarshalJSON(rawMessage, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	if req.Downloads == applingoapi.NoChange && req.Rating == applingoapi.NoChange {
		return openapi.DataResponseSuccess, nil
	}

	id := utils.GenerateDictionaryID(params.Name, params.Author)
	key := map[string]types.AttributeValue{
		applingodictionary.ColumnId:          &types.AttributeValueMemberS{Value: id},
		applingodictionary.ColumnSubcategory: &types.AttributeValueMemberS{Value: params.Subcategory},
	}

	updateBuilder := expression.UpdateBuilder{}

	switch req.Downloads {
	case applingoapi.Increase:
		updateBuilder = updateBuilder.Add(expression.Name(applingodictionary.ColumnDownloads), expression.Value(1))
	case applingoapi.Decrease:
		updateBuilder = updateBuilder.Add(expression.Name(applingodictionary.ColumnDownloads), expression.Value(-1))
	}

	switch req.Rating {
	case applingoapi.Increase:
		updateBuilder = updateBuilder.Add(expression.Name(applingodictionary.ColumnRating), expression.Value(1))
	case applingoapi.Decrease:
		updateBuilder = updateBuilder.Add(expression.Name(applingodictionary.ColumnRating), expression.Value(-1))
	}

	condition := expression.AttributeExists(expression.Name(applingodictionary.ColumnId))
	if err := dbDynamo.Update(ctx, applingodictionary.TableName, key, updateBuilder, condition); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return nil, &api.HandleError{Status: http.StatusNotFound, Err: errors.New("item not found")}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: errors.Wrap(err, "failed to update item")}
	}
	return openapi.DataResponseSuccess, nil
}
