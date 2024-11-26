package main

import (
	"context"
	"encoding/json"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/rs/zerolog"
)

func handleGet(_ context.Context, _ zerolog.Logger, _ json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	response := applingoapi.CategoriesData{
		FrontSide: []applingoapi.CategoryItemV1{
			{Code: "he"},
			{Code: "ru"},
			{Code: "en"},
		},
		BackSide: []applingoapi.CategoryItemV1{
			{Code: "he"},
			{Code: "ru"},
			{Code: "en"},
		},
	}
	return openapi.DataResponseCategories(response), nil
}
