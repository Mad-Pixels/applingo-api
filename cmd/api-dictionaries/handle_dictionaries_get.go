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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 60

func handleDictionariesGet(ctx context.Context, logger zerolog.Logger, log json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
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

	// logger.Info().
	// 	Any("subcategory", *params.Subcategory).
	// 	Any("level", *params.Level).
	// 	Any("public", *params.Public).
	// 	Msg("Income")

	queryInput, err := buildQueryInput(params, logger)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	logger.Info().Any("queryInput", queryInput).Msg("buildQueryInput")
	dynamoQueryInput, err := dbDynamo.BuildQueryInput(*queryInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	logger.Info().Any("dynamoQueryInput", dynamoQueryInput).Msg("Dynamo query")
	result, err := dbDynamo.Query(ctx, applingodictionary.TableName, dynamoQueryInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	// var (
	// 	wg      sync.WaitGroup
	// 	itemsCh = make(chan applingodictionary.SchemaItem, len(result.Items))
	// )
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

func buildQueryInput(params applingoapi.GetDictionariesV1Params, logger zerolog.Logger) (*cloud.QueryInput, error) {
	qb := applingodictionary.NewQueryBuilder()

	isPublic := true
	if params.Public != nil && !*params.Public {
		isPublic = false
	}
	sortBy := applingoapi.Date
	if params.SortBy != nil {
		sortBy = applingoapi.ParamDictionarySortEnum(*params.SortBy)
	}
	useRatingSort := sortBy == applingoapi.Rating

	logger.Info().Bool("ratingSort", useRatingSort).Any("sortBy", sortBy).Msg("sort param")
	var indexName string

	// Выбор индекса в зависимости от параметров запроса
	if params.Level != nil && params.Subcategory != nil {
		if useRatingSort {
			indexName = applingodictionary.IndexPublicLevelSubcategoryByRatingIndex

			qb.WithPublicLevelSubcategoryByRatingIndexHashKey(*params.Level, *params.Subcategory, applingodictionary.BoolToInt(isPublic))
			logger.Info().Msg("using WithPublicLevelSubcategoryByRatingIndexHashKey")
		} else {
			indexName = applingodictionary.IndexPublicLevelSubcategoryByDateIndex

			qb.WithPublicLevelSubcategoryByDateIndexHashKey(*params.Level, *params.Subcategory, applingodictionary.BoolToInt(isPublic))
			logger.Info().Msg("using WithPublicLevelSubcategoryByDateIndexHashKey")
		}
	} else if params.Level != nil {
		if useRatingSort {
			indexName = applingodictionary.IndexPublicLevelByRatingIndex

			qb.WithPublicLevelByRatingIndexHashKey(*params.Level, applingodictionary.BoolToInt(isPublic))
			logger.Info().Msg("using WithPublicLevelByRatingIndexHashKey")
		} else {
			qb.WithPublicLevelByDateIndexHashKey(*params.Level, applingodictionary.BoolToInt(isPublic))
			logger.Info().Msg("using WithPublicLevelByDateIndexHashKey")
		}
	} else if params.Subcategory != nil {
		if useRatingSort {
			indexName = applingodictionary.IndexPublicByRatingIndex // Используем общий индекс
			qb.WithPublicByRatingIndexHashKey(applingodictionary.BoolToInt(isPublic))
			// Добавляем subcategory как условие фильтра, а не ключа
			qb.WithSubcategory(*params.Subcategory)
			logger.Info().Msg("using WithPublicByRatingIndexHashKey with subcategory filter")
		} else {
			indexName = applingodictionary.IndexPublicByDateIndex // Используем общий индекс
			qb.WithPublicByDateIndexHashKey(applingodictionary.BoolToInt(isPublic))
			// Добавляем subcategory как условие фильтра, а не ключа
			qb.WithSubcategory(*params.Subcategory)
			logger.Info().Msg("using WithPublicByDateIndexHashKey with subcategory filter")
		}
	} else {
		if useRatingSort {
			indexName = applingodictionary.IndexPublicByRatingIndex

			qb.WithPublicByRatingIndexHashKey(applingodictionary.BoolToInt(isPublic))
			logger.Info().Msg("using WithPublicByRatingIndexHashKey")
		} else {
			indexName = applingodictionary.IndexPublicByDateIndex

			qb.WithPublicByDateIndexHashKey(applingodictionary.BoolToInt(isPublic))
			logger.Info().Msg("using WithPublicByDateIndexHashKey")
		}
	}

	// Всегда сортировать по убыванию для обоих типов сортировки
	qb.OrderByDesc()

	// Пагинация
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

	// Построение запроса
	_, keyCondition, filterCondition, exclusiveStartKey, err := qb.Build()
	if err != nil {
		return nil, err
	}
	logger.Info().Any("idx", indexName).Msg("dynamoIndex")

	// Создание QueryInput с явной настройкой ScanForward=false для сортировки от большего к меньшему
	queryInput := &cloud.QueryInput{
		IndexName:        indexName,
		KeyCondition:     keyCondition,
		ProjectionFields: applingodictionary.IndexProjections[indexName],
		Limit:            pageLimit,
		// Явно устанавливаем FALSE для сортировки от большего к меньшему
		ScanForward:       false,
		ExclusiveStartKey: exclusiveStartKey,
	}

	if filterCondition != nil {
		queryInput.FilterCondition = *filterCondition
	}

	return queryInput, nil
}
