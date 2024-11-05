package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/pkg/api"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"

	"github.com/rs/zerolog"
)

type handleGetRequest struct{}

type handleGetResponse struct {
	FrontCategory []categoryItem `json:"front_category"`
	BackCategory  []categoryItem `json:"back_category"`
}

type categoryItem struct {
	Name string `json:"name"`
}

func handleGet(_ context.Context, _ zerolog.Logger, data json.RawMessage) (any, *api.HandleError) {
	var req handleGetRequest
	if err := serializer.UnmarshalJSON(data, &req); err != nil {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}
	return handleGetResponse{
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
	}, nil
}
