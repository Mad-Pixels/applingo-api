package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func handleDelete(ctx context.Context, _ zerolog.Logger, raw json.RawMessage, pathParams map[string]string) (any, *api.HandleError) {
	id := pathParams["id"]
	if id == "" {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: errors.New("dictionary id is required")}
	}

	key := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: id},
	}

	// Сначала получаем словарь, чтобы узнать имя файла для удаления из S3
	result, err := dbDynamo.Get(ctx, applingodictionary.TableSchema.TableName, key)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if len(result.Item) == 0 {
		return nil, &api.HandleError{Status: http.StatusNotFound, Err: errors.New("dictionary not found")}
	}

	var item applingodictionary.SchemaItem
	if err = attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	// Удаляем файл из S3
	if err = s3Bucket.Delete(ctx, item.Filename, serviceDictionaryBucket); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	// Удаляем запись из DynamoDB
	if err = dbDynamo.Delete(ctx, applingodictionary.TableSchema.TableName, key); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	return nil, nil // DELETE возвращает 204 No Content
}
