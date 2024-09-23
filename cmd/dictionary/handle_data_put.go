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
	"github.com/Mad-Pixels/lingocards-api/pkg/tools"
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

func handleDataPut(ctx context.Context, _ zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataPutRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	item := make(map[string]types.AttributeValue)
	for _, attr := range gen_lingocards_dictionary.TableSchema.Attributes {
		switch attr.Name {
		case "id":
			item[attr.Name] = &types.AttributeValueMemberS{Value: tools.NewPersistentID(req.Author).UniqueID}
		case "name":
			item[attr.Name] = &types.AttributeValueMemberS{Value: req.Name}
		case "author":
			item[attr.Name] = &types.AttributeValueMemberS{Value: req.Author}
		case "category_main":
			item[attr.Name] = &types.AttributeValueMemberS{Value: req.CategoryMain}
		case "category_sub":
			item[attr.Name] = &types.AttributeValueMemberS{Value: req.CategorySub}
		case "description":
			item[attr.Name] = &types.AttributeValueMemberS{Value: req.Description}
		case "is_private":
			item[attr.Name] = &types.AttributeValueMemberN{Value: req.privateAttributeValue()}
		case "name_author":
			item[attr.Name] = &types.AttributeValueMemberS{Value: hex.EncodeToString(md5.New().Sum([]byte(req.Name + "-" + req.Author)))}
		}
	}
	if err := dbDynamo.Put(ctx, gen_lingocards_dictionary.TableSchema.TableName, item); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataPutResponse{
		Msg: "OK",
	}, nil
}
