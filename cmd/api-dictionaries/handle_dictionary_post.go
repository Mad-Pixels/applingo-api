package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
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

func handleDictionaryPost(ctx context.Context, _ zerolog.Logger, body json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if api.MustGetMetaData(ctx).IsDevice() || !api.MustGetMetaData(ctx).HasPermissions(auth.User) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	var req applingoapi.RequestPostDictionaryV1
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	levelSubcategoryIsPublic := fmt.Sprintf("%s#%s#%d", req.Level, req.Subcategory, applingodictionary.BoolToInt(req.Public))
	subcategoryIsPublic := fmt.Sprintf("%s#%d", req.Subcategory, applingodictionary.BoolToInt(req.Public))
	levelIsPublic := fmt.Sprintf("%s#%d", req.Level, applingodictionary.BoolToInt(req.Public))

	item := applingodictionary.SchemaItem{
		Id:          utils.GenerateDictionaryID(req.Name, req.Author),
		Name:        req.Name,
		Author:      req.Author,
		Category:    string(req.Category),
		Subcategory: req.Subcategory,
		Description: req.Description,
		IsPublic:    applingodictionary.BoolToInt(req.Public),
		Level:       req.Level,
		Created:     int(time.Now().Unix()),
		Rating:      0,

		// Composite keys
		LevelSubcategoryIsPublic: levelSubcategoryIsPublic,
		LevelIsPublic:            levelIsPublic,
		SubcategoryIsPublic:      subcategoryIsPublic,
	}
	dynamoItem, err := applingodictionary.PutItem(item)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	if err = dbDynamo.Put(
		ctx,
		applingodictionary.TableSchema.TableName,
		dynamoItem,
		expression.AttributeNotExists(expression.Name(applingodictionary.ColumnId)),
	); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return nil, &api.HandleError{Status: http.StatusConflict, Err: err}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return openapi.DataResponseSuccess, nil
}
