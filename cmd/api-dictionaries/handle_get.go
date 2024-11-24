package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/sort"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 40

type dictionaryItem struct {
	Name        string `json:"name" dynamodbav:"name"`
	Category    string `json:"category" dynamodbav:"category"`
	Subcategory string `json:"subcategory" dynamodbav:"subcategory"`
	Author      string `json:"author" dynamodbav:"author"`
	Dictionary  string `json:"dictionary" dynamodbav:"dictionary"`
	Description string `json:"description" dynamodbav:"description"`
	CreatedAt   int    `json:"created_at" dynamodbav:"created_at"`
	Rating      int    `json:"rating" dynamodbav:"rating"`
	IsPublic    int    `json:"is_public" dynamodbav:"is_public"`
}

type getDictionariesResponse struct {
	Items         []dictionaryItem `json:"items"`
	LastEvaluated string           `json:"last_evaluated,omitempty"`
}

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, queryParams map[string]string) (any, *api.HandleError) {
	// Преобразуем query parameters в структуру для построения запроса
	params := struct {
		SortBy        sort.QueryType
		Subcategory   string
		LastEvaluated string
		IsPublic      bool
	}{
		SortBy:        sort.QueryType(queryParams["sort_by"]),
		Subcategory:   queryParams["subcategory"],
		LastEvaluated: queryParams["last_evaluated"],
		IsPublic:      queryParams["is_public"] == "true",
	}

	queryInput, err := buildQueryInput(&params)
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
		itemsCh = make(chan dictionaryItem, len(result.Items))
	)
	response := getDictionariesResponse{
		Items: make([]dictionaryItem, 0, len(result.Items)),
	}

	for _, item := range result.Items {
		wg.Add(1)
		go func(item map[string]types.AttributeValue) {
			defer wg.Done()

			var dict dictionaryItem
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
		response.LastEvaluated = base64.StdEncoding.EncodeToString(lastEvaluatedKeyJSON)
	}

	return response, nil
}

func buildQueryInput(params *struct {
	SortBy        sort.QueryType
	Subcategory   string
	LastEvaluated string
	IsPublic      bool
}) (*cloud.QueryInput, error) {
	qb := applingodictionary.NewQueryBuilder()

	switch {
	case params.IsPublic && params.Subcategory != "":
		qb.WithSubcategory(params.Subcategory)
		qb.WithIsPublic(applingodictionary.BoolToInt(true))
	case params.IsPublic:
		qb.WithIsPublic(applingodictionary.BoolToInt(true))
	case params.Subcategory != "":
		qb.WithSubcategory(params.Subcategory)
	}
	qb.OrderByDesc()

	indexName, keyCondition, filterCondition, exclusiveStartKey, err := qb.Build()
	if err != nil {
		return nil, err
	}

	additionalFilter := expression.Name("dictionary").AttributeExists().And(
		expression.Name("dictionary").NotEqual(expression.Value("")),
	)

	var filterCond expression.ConditionBuilder
	if filterCondition != nil {
		filterCond = filterCondition.And(additionalFilter)
	} else {
		filterCond = additionalFilter
	}

	if params.LastEvaluated != "" {
		lastEvaluatedKeyJSON, err := base64.StdEncoding.DecodeString(params.LastEvaluated)
		if err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to decode base64")
		}
		var lastEvaluatedKeyMap map[string]interface{}
		if err := serializer.UnmarshalJSON(lastEvaluatedKeyJSON, &lastEvaluatedKeyMap); err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to unmarshal JSON")
		}
		exclusiveStartKey, err = attributevalue.MarshalMap(lastEvaluatedKeyMap)
		if err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to marshal attribute value")
		}
	}

	projectionFields := applingodictionary.IndexProjections[indexName]

	return &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCond,
		ProjectionFields:  projectionFields,
		Limit:             pageLimit,
		ScanForward:       false,
		ExclusiveStartKey: exclusiveStartKey,
	}, nil
}
