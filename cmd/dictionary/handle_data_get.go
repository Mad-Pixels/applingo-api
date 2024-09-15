package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"log"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

const pageLimit = 30

type handleDataGetRequest struct {
	CategoryMain  *string `json:"category_main,omitempty"`
	CategorySub   *string `json:"category_sub,omitempty"`
	Author        *string `json:"author,omitempty"`
	IsPrivate     *bool   `json:"is_private,omitempty"`
	IsPublish     *bool   `json:"is_publish,omitempty"`
	LastEvaluated *string `json:"last_evaluated,omitempty"`
}

type handleDataGetResponse struct {
	Items         []map[string]interface{} `json:"items"`
	LastEvaluated string                   `json:"last_evaluated,omitempty"`
}

func handleDataGet(ctx context.Context, logger zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	logger.Info().RawJSON("request_data", data).Msg("Received request")

	var req handleDataGetRequest
	if err := serializer.UnmarshalJSON(data, &req); err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal request")
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	logger.Info().Interface("parsed_request", req).Msg("Request parsed successfully")

	queryInput, err := buildQueryInput(&req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to build query input")
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	logger.Info().Interface("query_input", queryInput).Msg("Query input built successfully")

	dynamoQueryInput, err := dbDynamo.BuildQueryInput(*queryInput)
	if err != nil {
		logger.Error().Err(err).Interface("query_input", queryInput).Msg("Failed to build DynamoDB query input")
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	logger.Info().Interface("dynamo_query_input", dynamoQueryInput).Msg("DynamoDB query input built successfully")

	result, err := dbDynamo.Query(ctx, serviceDictionaryDynamo, dynamoQueryInput)
	if err != nil {
		logger.Error().Err(err).Str("table", serviceDictionaryDynamo).Interface("dynamo_query_input", dynamoQueryInput).Msg("Failed to execute DynamoDB query")
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	logger.Info().Int("item_count", len(result.Items)).Msg("Query executed successfully")

	response := handleDataGetResponse{
		Items: make([]map[string]interface{}, 0, len(result.Items)),
	}
	for _, item := range result.Items {
		mappedItem := make(map[string]interface{})
		if err := attributevalue.UnmarshalMap(item, &mappedItem); err != nil {
			logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
			continue
		}
		response.Items = append(response.Items, mappedItem)
	}
	if result.LastEvaluatedKey != nil {
		lastEvaluated, err := serializer.MarshalJSON(result.LastEvaluatedKey)
		if err == nil {
			response.LastEvaluated = string(lastEvaluated)
		}
	}
	return response, nil
}

func buildQueryInput(req *handleDataGetRequest) (*cloud.QueryInput, error) {
	var (
		keyCondition     expression.KeyConditionBuilder
		filterCondition  expression.ConditionBuilder
		indexName        string
		projectionFields = []string{"id", "name", "author", "category_main", "sub_category", "is_private", "is_publish"}
	)

	switch {
	case req.Author != nil && *req.Author != "":
		keyCondition = expression.Key("author").Equal(expression.Value(*req.Author))
		indexName = "AuthorIndex"
		log.Printf("Using AuthorIndex with author: %s", *req.Author)
	case req.CategoryMain != nil && *req.CategoryMain != "":
		keyCondition = expression.Key("category_main").Equal(expression.Value(*req.CategoryMain))
		indexName = "CategoryMainIndex"
		log.Printf("Using CategoryMainIndex with category_main: %s", *req.CategoryMain)
	case req.CategorySub != nil && *req.CategorySub != "":
		keyCondition = expression.Key("sub_category").Equal(expression.Value(*req.CategorySub))
		indexName = "CategorySubIndex"
		log.Printf("Using CategorySubIndex with sub_category: %s", *req.CategorySub)
	case req.IsPrivate != nil:
		keyCondition = expression.Key("is_private").Equal(expression.Value(boolToInt(*req.IsPrivate)))
		indexName = "IsPrivateIndex"
		log.Printf("Using IsPrivateIndex with is_private: %v", *req.IsPrivate)
	case req.IsPublish != nil:
		keyCondition = expression.Key("is_publish").Equal(expression.Value(boolToInt(*req.IsPublish)))
		indexName = "IsPublishIndex"
		log.Printf("Using IsPublishIndex with is_publish: %v", *req.IsPublish)
	default:
		return nil, errors.New("at least one query parameter is required")
	}

	if req.IsPrivate != nil && indexName != "IsPrivateIndex" {
		isPrivateValue := boolToInt(*req.IsPrivate)
		privateFilter := expression.Name("is_private").Equal(expression.Value(isPrivateValue))
		filterCondition = addFilterCondition(filterCondition, privateFilter)
		log.Printf("Added is_private filter condition: is_private = %d", isPrivateValue)
	}

	var exclusiveStartKey map[string]types.AttributeValue
	if req.LastEvaluated != nil && *req.LastEvaluated != "" {
		if err := json.Unmarshal([]byte(*req.LastEvaluated), &exclusiveStartKey); err != nil {
			return nil, errors.New("invalid last_evaluated key")
		}
	}

	queryInput := &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCondition,
		ProjectionFields:  projectionFields,
		Limit:             pageLimit,
		ScanForward:       true,
		ExclusiveStartKey: exclusiveStartKey,
	}
	log.Printf("Built QueryInput: %+v", queryInput)

	return queryInput, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func addFilterCondition(existingCondition expression.ConditionBuilder, newCondition expression.ConditionBuilder) expression.ConditionBuilder {
	if existingCondition.IsSet() {
		return existingCondition.And(newCondition)
	}
	return newCondition
}

func mapDynamoItemToResponse(item map[string]types.AttributeValue) (map[string]interface{}, error) {
	var mappedItem map[string]interface{}
	err := attributevalue.UnmarshalMap(item, &mappedItem)
	if err != nil {
		return nil, err
	}
	return mappedItem, nil
}
