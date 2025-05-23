package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingodictionary"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func handleDictionaryDelete(ctx context.Context, _ zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	if api.MustGetMetaData(ctx).IsDevice() || !api.MustGetMetaData(ctx).HasPermissions(auth.User) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	params := applingoapi.DeleteDictionaryV1Params{
		Name:        baseParams.GetStringDefault("name", ""),
		Author:      baseParams.GetStringDefault("author", ""),
		Subcategory: baseParams.GetStringDefault("subcategory", ""),
	}
	if err := validate.ValidateStruct(&params); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	id := utils.GenerateDictionaryID(params.Name, params.Author)
	result, err := dbDynamo.Get(ctx, applingodictionary.TableName, map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: id},
		"subcategory": &types.AttributeValueMemberS{Value: params.Subcategory},
	})
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: errors.Wrap(err, "failed to get item for deletion")}
	}
	if result.Item == nil {
		return nil, &api.HandleError{Status: http.StatusNotFound, Err: errors.New("item not found")}
	}

	if err := dbDynamo.Delete(ctx, applingodictionary.TableName, map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: id},
		"subcategory": &types.AttributeValueMemberS{Value: params.Subcategory},
	}); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: errors.Wrap(err, "failed to delete item")}
	}
	return nil, nil
}
