package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/data"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 2

type handleDataGetRequest struct {
	ID            string          `json:"id,omitempty"`
	Name          string          `json:"name,omitempty"`
	CategoryMain  string          `json:"category_main,omitempty"`
	CategorySub   string          `json:"category_sub,omitempty"`
	Author        string          `json:"author,omitempty"`
	IsPrivate     serializer.Bool `json:"is_private,omitempty"`
	IsPublish     serializer.Bool `json:"is_publish,omitempty"`
	LastEvaluated string          `json:"last_evaluated,omitempty"`
}

type handleDataGetResponse struct {
	Items         []map[string]interface{} `json:"items"`
	LastEvaluated string                   `json:"last_evaluated,omitempty"`
}

func handleDataGet(ctx context.Context, logger zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataGetRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	queryInput, err := buildQueryInput(&req)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	dynamoQueryInput, err := dbDynamo.BuildQueryInput(*queryInput)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	result, err := dbDynamo.Query(ctx, data.TableName, dynamoQueryInput)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	response := handleDataGetResponse{
		Items: make([]map[string]interface{}, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		var mappedItem map[string]interface{}

		if err = attributevalue.UnmarshalMap(item, &mappedItem); err != nil {
			logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
			continue
		}
		response.Items = append(response.Items, mappedItem)
	}
	if result.LastEvaluatedKey != nil {
		// Преобразуем LastEvaluatedKey в map[string]interface{}
		var lastEvaluatedKeyMap map[string]interface{}
		err := attributevalue.UnmarshalMap(result.LastEvaluatedKey, &lastEvaluatedKeyMap)
		if err != nil {
			return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
		}

		// Маршалим map в JSON
		lastEvaluatedKeyJSON, err := serializer.MarshalJSON(lastEvaluatedKeyMap)
		if err != nil {
			return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
		}

		// Кодируем в base64
		response.LastEvaluated = base64.StdEncoding.EncodeToString(lastEvaluatedKeyJSON)
	}
	return response, nil
}

func buildQueryInput(req *handleDataGetRequest) (*cloud.QueryInput, error) {
	qb := data.NewQueryBuilder()

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
	if req.IsPrivate.Set {
		qb.WithIsPrivate(boolToInt(req.IsPrivate.Value))
	}
	if req.IsPublish.Set {
		qb.WithIsPublish(boolToInt(req.IsPublish.Value))
	}

	indexName, keyCondition, filterCondition := qb.Build()

	if indexName == "" && !filterCondition.IsSet() {
		return nil, errors.New("at least one query parameter is required")
	}

	var exclusiveStartKey map[string]types.AttributeValue
	if req.LastEvaluated != "" {
		// Декодируем base64
		lastEvaluatedKeyJSON, err := base64.StdEncoding.DecodeString(req.LastEvaluated)
		if err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to decode base64")
		}

		// Демаршалим JSON в map[string]interface{}
		var lastEvaluatedKeyMap map[string]interface{}
		if err := serializer.UnmarshalJSON(lastEvaluatedKeyJSON, &lastEvaluatedKeyMap); err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to unmarshal JSON")
		}

		// Преобразуем map[string]interface{} обратно в map[string]types.AttributeValue
		exclusiveStartKey, err = attributevalue.MarshalMap(lastEvaluatedKeyMap)
		if err != nil {
			return nil, errors.New("invalid last_evaluated key: unable to marshal attribute value")
		}
	}

	return &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCondition,
		ProjectionFields:  data.IndexProjections[indexName],
		Limit:             pageLimit,
		ScanForward:       true,
		ExclusiveStartKey: exclusiveStartKey,
	}, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
