package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/Mad-Pixels/lingocards-api/pkg/api"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/dynamodb-interface/lingocardsdictionary"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

const (
	defaultRawCode = "0"
)

type handleDataPutRequest struct {
	Description  string `json:"description" validate:"required"`
	Dictionary   string `json:"dictionary" validate:"required"`
	Name         string `json:"name" validate:"required,min=4,max=32"`
	Author       string `json:"author" validate:"required"`
	Code         string `json:"code,omitempty"`
	CategoryMain string `json:"category_main" validate:"required"`
	CategorySub  string `json:"category_sub" validate:"required"`
}

type handleDataPutResponse struct {
	Msg string `json:"msg"`
}

func handleDataPut(ctx context.Context, _ zerolog.Logger, raw json.RawMessage) (any, *api.HandleError) {
	var req handleDataPutRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	id := hex.EncodeToString(md5.New().Sum([]byte(req.Name + "-" + req.Author)))

	schemaItem := lingocardsdictionary.SchemaItem{
		Id:           id,
		Name:         req.Name,
		Author:       req.Author,
		CategoryMain: req.CategoryMain,
		CategorySub:  req.CategorySub,
		Description:  req.Description,
		Code:         defaultRawCode,
		IsPublic:     1,
	}
	if req.Code != "" {
		schemaItem.Code = req.Code
		schemaItem.IsPublic = 0
	}

	item, err := lingocardsdictionary.PutItem(schemaItem)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := dbDynamo.Put(ctx, lingocardsdictionary.TableSchema.TableName, item, expression.AttributeNotExists(expression.Name("id"))); err != nil {
		var cfe *types.ConditionalCheckFailedException

		if errors.As(err, &cfe) {
			return nil, &api.HandleError{Status: http.StatusConflict, Err: errors.New("dictionary with id: '" + id + "' already exists")}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataPutResponse{
		Msg: "OK",
	}, nil
}
