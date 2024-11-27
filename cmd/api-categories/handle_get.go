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
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const pageLimit = 1000

func handleGet(ctx context.Context, logger zerolog.Logger, _ json.RawMessage, baseParams openapi.QueryParams) (any, *api.HandleError) {
	var paramSide *applingoapi.GetCategoriesV1ParamsSide
	if baseSide := baseParams.GetStringPtr("side"); baseSide != nil {
		switch *baseSide {
		case "front":
			convertedSide := applingoapi.Front
			paramSide = &convertedSide
		case "back":
			convertedSide := applingoapi.Back
			paramSide = &convertedSide
		default:
			return nil, &api.HandleError{
				Status: http.StatusBadRequest,
				Err:    errors.New("Invalid value for 'side'"),
			}
		}
	}
	params := applingoapi.GetCategoriesV1Params{
		Side: paramSide,
	}

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

	var (
		wg      sync.WaitGroup
		itemsCh = make(chan applingoapi.CategoryItemV1, len(result.Items))
	)
	response := applingoapi.CategoriesData{}
	for _, item := range result.Items {
		wg.Add(1)
		go func(item map[string]types.AttributeValue) {
			defer wg.Done()

			var category applingoapi.CategoryItemV1
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
		if item.Side != nil {
			switch *item.Side {
			case applingoapi.CategoryItemV1SideFront:
				item.Side = nil
				response.FrontSide = append(response.FrontSide, item)
			case applingoapi.CategoryItemV1SideBack:
				item.Side = nil
				response.BackSide = append(response.BackSide, item)
			}
		}
	}
	return openapi.DataResponseCategories(response), nil
}

func buildQueryInput(params applingoapi.GetCategoriesV1Params) (*cloud.QueryInput, error) {
	qb := applingosubcategory.NewQueryBuilder()
	if params.Side != nil {
		qb.WithSideIndexHashKey(string(*params.Side))
	} else {
		return &cloud.QueryInput{
			TableName:   applingosubcategory.TableName,
			ScanForward: false,
			Limit:       pageLimit,
		}, nil
	}
	qb.OrderByDesc()
	qb.Limit(pageLimit)

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
