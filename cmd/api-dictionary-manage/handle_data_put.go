package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

type handleDataPutRequest struct {
	Description string `json:"description" validate:"required"`
	Filename    string `json:"filename" validate:"required"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
	Author      string `json:"author" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Subcategory string `json:"subcategory" validate:"required"`
	IsPublic    bool   `json:"is_public" validate:"required"`
}

type handleDataPutResponse struct {
	Status string `json:"status"`
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

	schemaItem := applingodictionary.SchemaItem{
		Id:          id,
		Name:        req.Name,
		Author:      req.Author,
		Filename:    req.Filename,
		Category:    req.Category,
		Subcategory: req.Subcategory,
		Description: req.Description,
		IsPublic:    applingodictionary.BoolToInt(req.IsPublic),
		CreatedAt:   int(time.Now().Unix()),
		Rating:      0,
	}
	item, err := applingodictionary.PutItem(schemaItem)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if err = dbDynamo.Put(ctx, applingodictionary.TableSchema.TableName, item,
		expression.AttributeNotExists(expression.Name("id"))); err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			return nil, &api.HandleError{
				Status: http.StatusConflict,
				Err:    errors.New("dictionary with id: '" + id + "' already exists"),
			}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataPutResponse{Status: "OK"}, nil
}
