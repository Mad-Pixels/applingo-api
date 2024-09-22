package main

import (
	"context"
	"encoding/json"
	"github.com/Mad-Pixels/lingocards-api/data"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
)

const pageLimit = 20

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

	return handleDataGetResponse{
		Items: make([]map[string]interface{}, 0, len(result.Items)),
	}, nil
}

func buildQueryInput(req *handleDataGetRequest) (*cloud.QueryInput, error) {
	qb := data.NewQueryBuilder()
	switch {
	case req.Author != nil && *req.Author != "":
		qb.WithAuthor(*req.Author)
	case req.CategoryMain != nil && *req.CategoryMain != "":
		qb.WithCategoryMain(*req.CategoryMain)
	case req.CategorySub != nil && *req.CategorySub != "":
		qb.WithCategorySub(*req.CategorySub)
	case req.IsPrivate != nil:
		qb.WithIsPrivate(*req.IsPrivate)
	case req.IsPublish != nil:
		qb.WithIsPublish(*req.IsPublish)
	default:
		return nil, errors.New("at least one query parameter is required")
	}

	if req.IsPrivate != nil && qb.IndexName != data.IndexIsPrivateIndex {
		qb.WithIsPrivateFilter(*req.IsPrivate)
	}

	var exclusiveStartKey map[string]types.AttributeValue
	if req.LastEvaluated != nil && *req.LastEvaluated != "" {
		if err := json.Unmarshal([]byte(*req.LastEvaluated), &exclusiveStartKey); err != nil {
			return nil, errors.New("invalid last_evaluated key")
		}
	}

	return &cloud.QueryInput{
		IndexName:         qb.IndexName,
		KeyCondition:      qb.KeyCondition,
		FilterCondition:   qb.FilterCondition,
		ProjectionFields:  data.IndexProjections[qb.IndexName],
		Limit:             pageLimit,
		ScanForward:       true,
		ExclusiveStartKey: exclusiveStartKey,
	}, nil
}
