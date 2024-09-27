package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/data/gen_lingocards_dictionary"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
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

func handleDataPut(ctx context.Context, _ zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataPutRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	id := hex.EncodeToString(md5.New().Sum([]byte(req.Name + "-" + req.Author)))

	schemaItem := gen_lingocards_dictionary.SchemaItem{
		Id:                 id,
		Name:               req.Name,
		Author:             req.Author,
		CategoryMain:       req.CategoryMain,
		CategorySub:        req.CategorySub,
		Description:        req.Description,
		DictionaryFilename: req.Dictionary,
		Code:               defaultRawCode,
		IsPublic:           1,
	}
	if req.Code != "" {
		schemaItem.Code = req.Code
		schemaItem.IsPublic = 0
	}

	item, err := gen_lingocards_dictionary.PutItem(schemaItem)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := dbDynamo.Put(ctx, gen_lingocards_dictionary.TableSchema.TableName, item, expression.AttributeNotExists(expression.Name("id"))); err != nil {
		var cfe *types.ConditionalCheckFailedException

		if errors.As(err, &cfe) {
			return nil, &lambda.HandleError{Status: http.StatusConflict, Err: errors.New("dictionary with id: '" + id + "' already exists")}
		}
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataPutResponse{
		Msg: "OK",
	}, nil
}
