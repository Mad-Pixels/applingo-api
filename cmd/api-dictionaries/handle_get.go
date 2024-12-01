package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 60

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams, reqCtx api.ReqCtx) (any, *api.HandleError) {
	if !reqCtx.HasPermissions(auth.GetPermissionLevel(auth.Device)) {
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

	var (
		wg      sync.WaitGroup
		itemsCh = make(chan applingoapi.DictionaryItemV1, len(result.Items))
	)
	response := applingoapi.DictionariesData{
		Items: make([]applingoapi.DictionaryItemV1, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		wg.Add(1)
		go func(item map[string]types.AttributeValue) {
			defer wg.Done()

			var dict applingoapi.DictionaryItemV1
			if err := attributevalue.UnmarshalMap(item, &dict); err != nil {
				logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
				return
			}
			itemsCh <- dict
		}(item)
	}
	go func() {
		wg.Wait()
		close(itemsCh)
	}()

	for item := range itemsCh {
		response.Items = append(response.Items, item)
	}
	if result.LastEvaluatedKey != nil {
		var lastEvaluatedKeyMap map[string]interface{}
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

	isPublic := true
	if params.Public != nil && !*params.Public {
		isPublic = false
	}
	sortBy := applingoapi.Date
	if params.SortBy != nil {
		sortBy = applingoapi.ParamDictionariesSortEnum(*params.SortBy)
	}
	useRatingSort := sortBy == applingoapi.Rating

	if params.Level != nil && params.Subcategory != nil {
		if useRatingSort {
			qb.WithPublicLevelSubcategoryByRatingIndexHashKey(*params.Level, *params.Subcategory, applingodictionary.BoolToInt(isPublic))
		} else {
			qb.WithPublicLevelSubcategoryByDateIndexHashKey(*params.Level, *params.Subcategory, applingodictionary.BoolToInt(isPublic))
		}
	} else if params.Level != nil {
		if useRatingSort {
			qb.WithPublicLevelByRatingIndexHashKey(*params.Level, applingodictionary.BoolToInt(isPublic))
		} else {
			qb.WithPublicLevelByDateIndexHashKey(*params.Level, applingodictionary.BoolToInt(isPublic))
		}
	} else if params.Subcategory != nil {
		if useRatingSort {
			qb.WithPublicSubcategoryByRatingIndexHashKey(*params.Subcategory, applingodictionary.BoolToInt(isPublic))
		} else {
			qb.WithPublicSubcategoryByDateIndexHashKey(*params.Subcategory, applingodictionary.BoolToInt(isPublic))
		}
	} else {
		if useRatingSort {
			qb.WithPublicByRatingIndexHashKey(applingodictionary.BoolToInt(isPublic))
		} else {
			qb.WithPublicByDateIndexHashKey(applingodictionary.BoolToInt(isPublic))
		}
	}
	qb.OrderByDesc()

	if params.LastEvaluated != nil {
		lastEvaluatedKeyJSON, err := base64.StdEncoding.DecodeString(*params.LastEvaluated)
		if err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to decode base64")
		}
		var lastEvaluatedKeyMap map[string]types.AttributeValue
		if err := json.Unmarshal(lastEvaluatedKeyJSON, &lastEvaluatedKeyMap); err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to unmarshal JSON")
		}
		qb.StartFrom(lastEvaluatedKeyMap)
	}
	qb.Limit(pageLimit)

	additionalFilter := expression.Name("dictionary").AttributeExists().And(
		expression.Name("dictionary").NotEqual(expression.Value("")),
	)
	indexName, keyCondition, filterCondition, exclusiveStartKey, err := qb.Build()
	if err != nil {
		return nil, err
	}

	var filterCond expression.ConditionBuilder
	if filterCondition != nil {
		filterCond = filterCondition.And(additionalFilter)
	} else {
		filterCond = additionalFilter
	}
	return &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCond,
		ProjectionFields:  applingodictionary.IndexProjections[indexName],
		Limit:             pageLimit,
		ScanForward:       false,
		ExclusiveStartKey: exclusiveStartKey,
	}, nil
}
