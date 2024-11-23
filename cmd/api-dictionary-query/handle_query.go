package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/pkg/sort"

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

type handleDataQueryRequest struct {
	SortBy        sort.QueryType `json:"sort_by,omitempty"`
	Subcategory   string         `json:"subcategory,omitempty"`
	LastEvaluated string         `json:"last_evaluated,omitempty"`
	IsPublic      bool           `json:"is_public,omitempty"`
}

type handleDataQueryResponse struct {
	Items         []dataQueryItem `json:"items"`
	LastEvaluated string          `json:"last_evaluated,omitempty"`
}

type dataQueryItem struct {
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

func handleDataQuery(ctx context.Context, logger zerolog.Logger, raw json.RawMessage) (any, *api.HandleError) {
	var req handleDataQueryRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	queryInput, err := buildQueryInput(&req)
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
		itemsCh = make(chan dataQueryItem, len(result.Items))
	)
	response := handleDataQueryResponse{
		Items: make([]dataQueryItem, 0, len(result.Items)),
	}
	for _, dynamoItem := range result.Items {
		wg.Add(1)
		go func(dynamoItem map[string]types.AttributeValue) {
			defer wg.Done()

			var item dataQueryItem
			if err := attributevalue.UnmarshalMap(dynamoItem, &item); err != nil {
				logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
				return
			}
			itemsCh <- item
		}(dynamoItem)
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

func buildQueryInput(req *handleDataQueryRequest) (*cloud.QueryInput, error) {
	qb := applingodictionary.NewQueryBuilder()

	switch {
	case req.IsPublic && req.Subcategory != "":
		qb.WithSubcategory(req.Subcategory)
		qb.WithIsPublic(applingodictionary.BoolToInt(true))
	case req.IsPublic:
		qb.WithIsPublic(applingodictionary.BoolToInt(true))
	case req.Subcategory != "":
		qb.WithSubcategory(req.Subcategory)
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
	if req.LastEvaluated != "" {
		lastEvaluatedKeyJSON, err := base64.StdEncoding.DecodeString(req.LastEvaluated)
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
