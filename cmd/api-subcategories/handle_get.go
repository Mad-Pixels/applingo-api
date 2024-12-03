package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingosubcategory"
	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 1000

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).HasPermissions(auth.Device) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	validSideValues := map[applingoapi.BaseSideEnum]struct{}{
		applingoapi.Front: {},
		applingoapi.Back:  {},
	}
	paramSide, err := openapi.ParseEnumParam(baseParams.GetStringPtr("side"), validSideValues)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: errors.Wrap(err, "invalid value for 'side' param")}
	}
	params := applingoapi.GetSubcategoriesV1Params{
		Side: paramSide,
	}
	if err := validate.ValidateStruct(&params); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	var items []map[string]types.AttributeValue
	if params.Side != nil {
		queryInput, err := buildQueryInput(params)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
		}

		dynamoQueryInput, err := dbDynamo.BuildQueryInput(*queryInput)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		result, err := dbDynamo.Query(ctx, applingosubcategory.TableName, dynamoQueryInput)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		items = result.Items
	} else {
		scanInput := dbDynamo.BuildScanInput(applingosubcategory.TableName, pageLimit, nil)
		result, err := dbDynamo.Scan(ctx, applingosubcategory.TableName, scanInput)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		items = result.Items
	}

	var (
		wg       sync.WaitGroup
		itemsCh  = make(chan applingosubcategory.SchemaItem, len(items))
		response = applingoapi.CategoriesData{}
	)
	for _, item := range items {
		wg.Add(1)
		go func(item map[string]types.AttributeValue) {
			defer wg.Done()

			var category applingosubcategory.SchemaItem
			if err := attributevalue.UnmarshalMap(item, &category); err != nil {
				logger.Warn().Err(err).Msg("Failed to unmarshal DynamoDB item")
				return
			}
			itemsCh <- category
		}(item)
	}
	go func() {
		wg.Wait()
		close(itemsCh)
	}()

	for item := range itemsCh {
		switch applingoapi.BaseSideEnum(item.Side) {
		case applingoapi.Front:
			response.FrontSide = append(response.FrontSide, applingoapi.SubcategoryItemV1{
				Code: item.Code,
			})
		case applingoapi.Back:
			response.BackSide = append(response.BackSide, applingoapi.SubcategoryItemV1{
				Code: item.Code,
			})
		}
	}
	return openapi.DataResponseSubcategories(response), nil
}

func buildQueryInput(params applingoapi.GetSubcategoriesV1Params) (*cloud.QueryInput, error) {
	qb := applingosubcategory.NewQueryBuilder()
	if params.Side != nil {
		qb.WithSideIndexHashKey(string(*params.Side))
		qb.Limit(pageLimit)
		qb.OrderByDesc()

		indexName, keyCondition, _, exclusiveStartKey, err := qb.Build()
		if err != nil {
			return nil, err
		}
		return &cloud.QueryInput{
			IndexName:         indexName,
			KeyCondition:      keyCondition,
			ProjectionFields:  applingosubcategory.IndexProjections[indexName],
			Limit:             pageLimit,
			ScanForward:       false,
			ExclusiveStartKey: exclusiveStartKey,
		}, nil
	}
	return nil, nil
}
