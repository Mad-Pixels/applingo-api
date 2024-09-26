package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/data/gen_lingocards_dictionary"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 20

type handleDataQueryRequest struct {
	ID            string          `json:"id,omitempty"`
	Name          string          `json:"name,omitempty"`
	CategoryMain  string          `json:"category_main,omitempty"`
	CategorySub   string          `json:"category_sub,omitempty"`
	Author        string          `json:"author,omitempty"`
	IsPrivate     serializer.Bool `json:"is_private,omitempty"`
	IsPublish     serializer.Bool `json:"is_publish,omitempty"`
	LastEvaluated string          `json:"last_evaluated,omitempty"`
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
	DictionaryUrl string `json:"dictionary_url,omitempty" dynamodbav:"dictionary_url"`
	Description   string `json:"description,omitempty" dynamodbav:"description"`
}

func handleDataQuery(ctx context.Context, logger zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataQueryRequest
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
	result, err := dbDynamo.Query(ctx, gen_lingocards_dictionary.TableName, dynamoQueryInput)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	response := handleDataQueryResponse{
		Items: make([]dataQueryItem, 0, len(result.Items)),
	}
	for _, dynamoItem := range result.Items {
		var item dataQueryItem

		logger.Info().Interface("DynamoDBItem", dynamoItem).Msg("DynamoDB item")

		if err = attributevalue.UnmarshalMap(dynamoItem, &item); err != nil {
			logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
			continue
		}
		response.Items = append(response.Items, item)
	}
	if result.LastEvaluatedKey != nil {
		var lastEvaluatedKeyMap map[string]interface{}

		if err = attributevalue.UnmarshalMap(result.LastEvaluatedKey, &lastEvaluatedKeyMap); err != nil {
			return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		lastEvaluatedKeyJSON, err := serializer.MarshalJSON(lastEvaluatedKeyMap)
		if err != nil {
			return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		response.LastEvaluated = base64.StdEncoding.EncodeToString(lastEvaluatedKeyJSON)
	}
	return response, nil
}

func buildQueryInput(req *handleDataQueryRequest) (*cloud.QueryInput, error) {
	qb := gen_lingocards_dictionary.NewQueryBuilder()

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
		qb.WithIsPrivate(gen_lingocards_dictionary.BoolToInt(req.IsPrivate.Value))
	}
	if req.IsPublish.Set {
		qb.WithIsPublish(gen_lingocards_dictionary.BoolToInt(req.IsPublish.Value))
	}

	indexName, keyCondition, filterCondition := qb.Build()

	if indexName == "" && !filterCondition.IsSet() {
		return nil, errors.New("at least one query parameter is required")
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

	return &cloud.QueryInput{
		IndexName:         indexName,
		KeyCondition:      keyCondition,
		FilterCondition:   filterCondition,
		ProjectionFields:  gen_lingocards_dictionary.IndexProjections[indexName],
		Limit:             pageLimit,
		ScanForward:       true,
		ExclusiveStartKey: exclusiveStartKey,
	}, nil
}
