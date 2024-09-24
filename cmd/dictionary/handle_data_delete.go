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
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

type handleDataDeleteRequest struct {
	Name   string `json:"name" validate:"required,min=4,max=32"`
	Author string `json:"author" validate:"required"`
}

type handleDataDeleteResponse struct {
	Msg string `json:"msg"`
}

func handleDataDelete(ctx context.Context, logger zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataDeleteRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	var (
		id  = hex.EncodeToString(md5.New().Sum([]byte(req.Name + "-" + req.Author)))
		key = map[string]types.AttributeValue{
			"id":   &types.AttributeValueMemberS{Value: id},
			"name": &types.AttributeValueMemberS{Value: req.Name},
		}
	)
	result, err := dbDynamo.Get(ctx, gen_lingocards_dictionary.TableSchema.TableName, key)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if len(result.Item) == 0 {
		return nil, &lambda.HandleError{Status: http.StatusNotFound, Err: err}
	}

	var item gen_lingocards_dictionary.SchemaItem
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if err := s3Bucket.Delete(ctx, item.DictionaryFilename, serviceDictionaryBucket); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	if err := dbDynamo.Delete(ctx, gen_lingocards_dictionary.TableSchema.TableName, key); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataDeleteResponse{
		Msg: "OK",
	}, nil
}
