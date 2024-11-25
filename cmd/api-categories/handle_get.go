package main

import (
	"context"
	"encoding/json"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/v1/categories"
	"github.com/Mad-Pixels/applingo-api/pkg/api"

	"github.com/rs/zerolog"
)

func handleGet(_ context.Context, _ zerolog.Logger, _ json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	response := categories.GetResponse{
		FrontCategory: []categories.Item{
			{Name: "he"},
			{Name: "ru"},
			{Name: "en"},
		},
		BackCategory: []categories.Item{
			{Name: "he"},
			{Name: "ru"},
			{Name: "en"},
		},
	}
	return response, nil
}
