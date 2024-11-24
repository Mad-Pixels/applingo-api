package main

import (
	"context"
	"encoding/json"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/rs/zerolog"
)

type handleQueryResponse struct {
	FrontCategory []categoryItem `json:"front_category"`
	BackCategory  []categoryItem `json:"back_category"`
}

type categoryItem struct {
	Name string `json:"name"`
}

func handleQuery(_ context.Context, _ zerolog.Logger, data json.RawMessage, queryParams map[string]string) (any, *api.HandleError) {
	response := handleQueryResponse{
		FrontCategory: []categoryItem{
			{Name: "he"},
			{Name: "ru"},
			{Name: "en"},
		},
		BackCategory: []categoryItem{
			{Name: "he"},
			{Name: "ru"},
			{Name: "en"},
		},
	}
	return response, nil
}
