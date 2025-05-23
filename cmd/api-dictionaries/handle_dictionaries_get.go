package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 150

func handleDictionariesGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	validSortValues := map[applingoapi.BaseDictSortEnum]struct{}{
		applingoapi.Date:   {},
		applingoapi.Rating: {},
	}
	paramSort, err := openapi.ParseEnumParam(baseParams.GetStringPtr("sort_by"), validSortValues)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: errors.Wrap(err, "invalid value for 'sort_by' param")}
	}
	params := applingoapi.GetDictionariesV1Params{
		Subcategory:   baseParams.GetStringPtr("subcategory"),
		LastEvaluated: baseParams.GetStringPtr("last_evaluated"),
		Level:         baseParams.GetStringPtr("level"),
		Public:        baseParams.GetBoolPtr("public"),
		SortBy:        paramSort,
	}
	if err := validate.ValidateStruct(&params); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	queryInput, err := buildQueryInput(params)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	dynamoQueryInput, err := dbDynamo.BuildQueryInput(*queryInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	result, err := dbDynamo.Query(ctx, applingodictionary.TableName, dynamoQueryInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	response := applingoapi.DictionariesData{
		Items: make([]applingoapi.DictionaryItemV1, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		var dict applingodictionary.SchemaItem
		if err := attributevalue.UnmarshalMap(item, &dict); err != nil {
			logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
			continue
		}

		response.Items = append(response.Items, applingoapi.DictionaryItemV1{
			Id:          dict.Id,
			Category:    applingoapi.BaseCategoryEnum(dict.Category),
			Public:      applingodictionary.IntToBool(dict.IsPublic),
			Dictionary:  utils.RecordToFileID(dict.Id),
			Downloads:   int64(dict.Downloads),
			Created:     int64(dict.Created),
			Rating:      int32(dict.Rating),
			Words:       int32(dict.Words),
			Subcategory: dict.Subcategory,
			Description: dict.Description,
			Author:      dict.Author,
			Name:        dict.Name,
			Level:       dict.Level,
			Topic:       dict.Topic,
		})
	}
	if result.LastEvaluatedKey != nil {
		var lastEvaluatedKeyMap map[string]any
		if err = attributevalue.UnmarshalMap(result.LastEvaluatedKey, &lastEvaluatedKeyMap); err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		lastEvaluatedKeyJSON, err := serializer.MarshalJSON(lastEvaluatedKeyMap)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		page := base64.StdEncoding.EncodeToString(lastEvaluatedKeyJSON)
		response.LastEvaluated = &page
	}
	return openapi.DataResponseDictionaries(response), nil
}

func buildQueryInput(params applingoapi.GetDictionariesV1Params) (*cloud.QueryInput, error) {
	qb := applingodictionary.NewQueryBuilder()

	isPublic := true //nolint:staticcheck
	if params.Public != nil && !*params.Public {
		isPublic = false
	}
	qb.WithIsPublic(applingodictionary.BoolToInt(isPublic))

	if params.Level != nil {
		qb.WithLevel(*params.Level)
	}
	if params.Subcategory != nil {
		qb.WithSubcategory(*params.Subcategory)
	}

	sortBy := applingoapi.Date
	if params.SortBy != nil {
		sortBy = applingoapi.ParamDictionarySortEnum(*params.SortBy)
	}
	useRatingSort := sortBy == applingoapi.Rating
	if useRatingSort {
		qb.WithPreferredSortKey(string(applingoapi.Rating))
		qb.WithRatingGreaterThan(-1)
	} else {
		qb.WithPreferredSortKey(string(applingoapi.Date))
		qb.WithCreatedGreaterThan(0)
	}

	if params.LastEvaluated != nil {
		lastEvaluatedKeyJSON, err := base64.StdEncoding.DecodeString(*params.LastEvaluated)
		if err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to decode base64")
		}

		var jsonMap map[string]any
		if err := json.Unmarshal(lastEvaluatedKeyJSON, &jsonMap); err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to unmarshal JSON")
		}
		lastEvaluatedKey, err := applingodictionary.ConvertMapToAttributeValues(jsonMap)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert last_evaluated to DynamoDB types")
		}

		qb.StartFrom(lastEvaluatedKey)
	}
	qb.Limit(pageLimit)

	indexName, keyCondition, filterCondition, exclusiveStartKey, err := qb.Build()
	if err != nil {
		return nil, err
	}
	queryInput := &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		ProjectionFields:  applingodictionary.IndexProjections[indexName],
		Limit:             pageLimit,
		ScanForward:       false,
		ExclusiveStartKey: exclusiveStartKey,
	}
	if filterCondition != nil {
		queryInput.FilterCondition = *filterCondition
	}
	return queryInput, nil
}
