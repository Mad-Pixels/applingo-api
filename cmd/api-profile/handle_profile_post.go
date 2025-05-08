package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprofile"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func handleProfilePost(ctx context.Context, _ zerolog.Logger, body json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).IsDevice() {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	var req applingoapi.RequestPostProfileV1
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	item := applingoprofile.SchemaItem{
		Id:    req.Id,
		Level: 1,
		Xp:    0,
	}
	dynamoItem, err := applingoprofile.PutItem(item)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	if err := dbDynamo.Put(
		ctx,
		applingoprofile.TableSchema.TableName,
		dynamoItem,
		expression.AttributeNotExists(expression.Name(applingoprofile.ColumnId)),
	); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return nil, &api.HandleError{Status: http.StatusConflict, Err: err}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return openapi.DataResponseSuccess, nil
}
