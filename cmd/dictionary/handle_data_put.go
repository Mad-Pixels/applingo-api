package main

import (
	"context"
	"encoding/json"
	"github.com/Mad-Pixels/lingocards-api/pkg/tools"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

type handleDataPutRequest struct {
	Description  string `json:"description" validate:"required"`
	Dictionary   string `json:"dictionary" validate:"required"`
	Name         string `json:"name" validate:"required,min=4,max=32"`
	Author       string `json:"author" validate:"required"`
	CategoryMain string `json:"category_main" validate:"required"`
	CategorySub  string `json:"category_sub" validate:"required"`
	Private      bool   `json:"private"`
}

func (r handleDataPutRequest) privateAttributeValue() string {
	if r.Private {
		return "0"
	}
	return "1"
}

type handleDataPutResponse struct {
	Msg string `json:"msg"`
}

func handleDataPut(ctx context.Context, _ zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataPutRequest
	if err := serializer.UnmarshalJSON(data, &req); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	item := map[string]types.AttributeValue{
		"id":            &types.AttributeValueMemberS{Value: tools.NewPersistentID(req.Author).UniqueID},
		"name":          &types.AttributeValueMemberS{Value: req.Name},
		"author":        &types.AttributeValueMemberS{Value: req.Author},
		"category_main": &types.AttributeValueMemberS{Value: req.CategoryMain},
		"category_sub":  &types.AttributeValueMemberS{Value: req.CategorySub},
		"description":   &types.AttributeValueMemberS{Value: req.Description},
		"is_private":    &types.AttributeValueMemberN{Value: req.privateAttributeValue()},
	}
	if err := dbDynamo.Put(ctx, serviceDictionaryDynamo, item); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}
	return handleDataPutResponse{
		Msg: "OK",
	}, nil
}
