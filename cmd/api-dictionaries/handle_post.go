package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func handlePost(ctx context.Context, logger zerolog.Logger, body json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	var req applingoapi.RequestPostDictionariesV1
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
		Id:          generateDictionaryID(req.Name, req.Author),
		Name:        req.Name,
		Author:      req.Author,
		Filename:    req.Filename,
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
		expression.AttributeNotExists(expression.Name("id")),
	); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return nil, &api.HandleError{Status: http.StatusConflict, Err: err}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return openapi.DataResponseSuccess, nil
}

func generateDictionaryID(name, author string) string {
	hash := md5.New()
	hash.Write([]byte(name + "-" + author))
	return hex.EncodeToString(hash.Sum(nil))
}
