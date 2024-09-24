package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/data/gen_lingocards_dictionary"
	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
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

type handleDataPutResponse struct {
	Msg string `json:"msg"`
}

func handleDataPut(ctx context.Context, logger zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataPutRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	id := hex.EncodeToString(md5.New().Sum([]byte(req.Name + "-" + req.Author)))

	exists, err := checkIfDictionaryExist(ctx, id, req.Name)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if exists {
		return nil, &lambda.HandleError{Status: http.StatusConflict, Err: errors.New("dictionary with id: '" + id + "' already exist")}
	}

	item, err := gen_lingocards_dictionary.PutItem(gen_lingocards_dictionary.SchemaItem{
		IsPrivate:    gen_lingocards_dictionary.BoolToInt(req.Private),
		Id:           id,
		Name:         req.Name,
		Author:       req.Author,
		CategoryMain: req.CategoryMain,
		CategorySub:  req.CategorySub,
		Description:  req.Description,
		Dictionary:   req.Dictionary,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Error creating item for DynamoDB")
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := dbDynamo.Put(ctx, gen_lingocards_dictionary.TableSchema.TableName, item); err != nil {
		logger.Error().Err(err).Msg("Error inserting item into DynamoDB")
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataPutResponse{
		Msg: "OK",
	}, nil
}

func checkIfDictionaryExist(ctx context.Context, id, name string) (bool, error) {
	result, err := dbDynamo.Get(
		ctx,
		gen_lingocards_dictionary.TableSchema.TableName,
		map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: id},
			"name": &types.AttributeValueMemberS{Value: name},
		},
	)
	if err != nil {
		return false, err
	}
	return len(result.Item) > 0, nil
}
