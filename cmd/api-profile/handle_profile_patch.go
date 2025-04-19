package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprofile"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rs/zerolog"
)

func handleProfilePatch(ctx context.Context, _ zerolog.Logger, body json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	var req applingoapi.RequestPatchProfileV1
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	key := map[string]types.AttributeValue{
		applingoprofile.ColumnId: &types.AttributeValueMemberS{Value: req.Id},
	}
	updateBuilder := expression.UpdateBuilder{}.
		Set(expression.Name(applingoprofile.ColumnLevel), expression.Value(req.Level)).
		Set(expression.Name(applingoprofile.ColumnXp), expression.Value(req.Xp))

	condition := expression.AttributeExists(expression.Name(applingoprofile.ColumnId))
	if err := dbDynamo.Update(
		ctx,
		applingoprofile.TableSchema.TableName,
		key,
		updateBuilder,
		condition,
	); err != nil {
		var conditionErr *types.ConditionalCheckFailedException
		if errors.As(err, &conditionErr) {
			return nil, &api.HandleError{Status: http.StatusNotFound, Err: errors.New("item not found")}
		}
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return nil, nil
}
