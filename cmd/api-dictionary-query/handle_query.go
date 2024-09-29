package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/Mad-Pixels/lingocards-api/dynamodb-interface/lingocardsdictionary"
	"github.com/Mad-Pixels/lingocards-api/pkg/api"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 20

type handleDataQueryRequest struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	CategoryMain  string `json:"category_main,omitempty"`
	CategorySub   string `json:"category_sub,omitempty"`
	Author        string `json:"author,omitempty"`
	IsPublic      bool   `json:"is_public,omitempty"`
	Code          string `json:"code,omitempty"`
	LastEvaluated string `json:"last_evaluated,omitempty"`
}

type handleDataQueryResponse struct {
	Items         []dataQueryItem `json:"items"`
	LastEvaluated string          `json:"last_evaluated,omitempty"`
}

type dataQueryItem struct {
	Name          string `json:"name,omitempty" dynamodbav:"name"`
	CategoryMain  string `json:"category_main,omitempty" dynamodbav:"category_main"`
	CategorySub   string `json:"category_sub,omitempty" dynamodbav:"category_sub"`
	Author        string `json:"author,omitempty" dynamodbav:"author"`
	DictionaryKey string `json:"dictionary_key,omitempty" dynamodbav:"dictionary_key"`
	Description   string `json:"description,omitempty" dynamodbav:"description"`
}

func handleDataQuery(ctx context.Context, logger zerolog.Logger, raw json.RawMessage) (any, *api.HandleError) {
	var req handleDataQueryRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if req.Code == "" {
		req.IsPublic = true
	}

	queryInput, err := buildQueryInput(&req)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	dynamoQueryInput, err := dbDynamo.BuildQueryInput(*queryInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	result, err := dbDynamo.Query(ctx, lingocardsdictionary.TableName, dynamoQueryInput)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	response := handleDataQueryResponse{
		Items: make([]dataQueryItem, 0, len(result.Items)),
	}

	var wg sync.WaitGroup
	itemsCh := make(chan dataQueryItem, len(result.Items))
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
	qb := lingocardsdictionary.NewQueryBuilder()

	if req.ID != "" {
		qb.WithId(req.ID)
	}
	if req.Name != "" {
		qb.WithName(req.Name)
	}
	if req.Author != "" {
		qb.WithAuthor(req.Author)
	}
	if req.CategoryMain != "" {
		qb.WithCategoryMain(req.CategoryMain)
	}
	if req.CategorySub != "" {
		qb.WithCategorySub(req.CategorySub)
	}
	if req.Code != "" {
		qb.WithCode(req.Code)
	}
	if req.IsPublic {
		qb.WithIsPublic(lingocardsdictionary.BoolToInt(req.IsPublic))
	}
	indexName, keyCondition, filterCondition, err := qb.Build()
	if err != nil {
		return nil, err
	}

	additionalFilter := expression.Name("dictionary_key").AttributeExists().And(
		expression.Name("dictionary_key").NotEqual(expression.Value("")),
	)

	var filterCond expression.ConditionBuilder
	if filterCondition != nil {
		combinedFilter := filterCondition.And(additionalFilter)
		filterCond = combinedFilter
	} else {
		filterCond = additionalFilter
	}

	var exclusiveStartKey map[string]types.AttributeValue
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
	projectionFields := lingocardsdictionary.IndexProjections[indexName]

	return &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCond,
		ProjectionFields:  projectionFields,
		Limit:             pageLimit,
		ScanForward:       true,
		ExclusiveStartKey: exclusiveStartKey,
	}, nil
}
