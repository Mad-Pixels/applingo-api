package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 40

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	params := applingoapi.GetDictionariesV1Params{
		SortBy:        baseParams.GetStringPtr("sort_by"),
		Subcategory:   baseParams.GetStringPtr("subcategory"),
		LastEvaluated: baseParams.GetStringPtr("last_evaluated"),
		Level:         baseParams.GetStringPtr("level"),
		Public:        baseParams.GetBoolPtr("public"),
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

	// Устанавливаем сортировку по умолчанию - "date"
	sortBy := "date"
	if params.SortBy != nil && (*params.SortBy == "date" || *params.SortBy == "rating") {
		sortBy = *params.SortBy
	}

	// Выбираем индекс в зависимости от предоставленных параметров
	if params.Level != nil && params.Subcategory != nil {
		compositeKey := fmt.Sprintf("public_%s_%s", *params.Level, *params.Subcategory)
		if sortBy == "date" {
			qb.WithPublicLevelSubcategoryByDateIndexHashKey(compositeKey)
		} else {
			qb.WithPublicLevelSubcategoryByRatingIndexHashKey(compositeKey)
		}
	} else if params.Level != nil {
		compositeKey := fmt.Sprintf("public_%s", *params.Level)
		if sortBy == "date" {
			qb.WithPublicLevelByDateIndexHashKey(compositeKey)
		} else {
			qb.WithPublicLevelByRatingIndexHashKey(compositeKey)
		}
	} else if params.Subcategory != nil {
		compositeKey := fmt.Sprintf("public_%s", *params.Subcategory)
		if sortBy == "date" {
			qb.WithPublicSubcategoryByDateIndexHashKey(compositeKey)
		} else {
			qb.WithPublicSubcategoryByRatingIndexHashKey(compositeKey)
		}
	} else {
		return nil, errors.New("Level or Subcategory must be provided when querying public dictionaries")
	}

	// Устанавливаем сортировку по убыванию
	qb.OrderByDesc()

	// Обрабатываем LastEvaluatedKey для пагинации
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

	// Устанавливаем лимит на количество элементов
	qb.Limit(pageLimit)

	// Создаем дополнительный фильтр для проверки наличия атрибута "dictionary"
	additionalFilter := expression.Name("dictionary").AttributeExists().And(
		expression.Name("dictionary").NotEqual(expression.Value("")),
	)

	// Строим запрос
	indexName, keyCondition, filterCondition, _, err := qb.Build()
	if err != nil {
		return nil, err
	}

	// Комбинируем существующий фильтр с дополнительным
	var filterCond expression.ConditionBuilder
	if filterCondition != nil {
		filterCond = filterCondition.And(additionalFilter)
	} else {
		filterCond = additionalFilter
	}

	// Получаем поля для проекции
	projectionFields := applingodictionary.IndexProjections[indexName]

	// Возвращаем сформированный запрос
	return &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCond,
		ProjectionFields:  projectionFields,
		Limit:             pageLimit,
		ScanForward:       false,
		ExclusiveStartKey: qb.ExclusiveStartKey,
	}, nil
}
