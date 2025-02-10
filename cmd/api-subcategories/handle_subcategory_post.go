package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingosubcategory"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func handleSubcategoryPost(ctx context.Context, _ zerolog.Logger, body json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if api.MustGetMetaData(ctx).IsDevice() || !api.MustGetMetaData(ctx).HasPermissions(auth.User) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	var req applingoapi.RequestPostSubcategoryV1
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	item := applingosubcategory.SchemaItem{
		Id:          utils.GenerateSubcategoryID(req.Code, string(req.Side)),
		Code:        req.Code,
		Side:        string(req.Side),
		Description: req.Description,
	}
	dynamoItem, err := applingosubcategory.PutItem(item)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if err = dbDynamo.Put(
		ctx,
		applingosubcategory.TableName,
		dynamoItem,
		expression.AttributeNotExists(expression.Name(applingosubcategory.ColumnId)),
	); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return nil, &api.HandleError{Status: http.StatusConflict, Err: err}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return openapi.DataResponseSuccess, nil
}
