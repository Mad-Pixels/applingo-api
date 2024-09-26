package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/rs/zerolog"
)

type handleGetRequest struct{}

type handleGetResponse struct {
	Categories []category `json:"categories"`
}

type category struct {
	Name          string   `json:"name"`
	SubCategories []string `json:"sub_categories"`
}

func handleGet(_ context.Context, _ zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	var req handleGetRequest
	if err := serializer.UnmarshalJSON(data, &req); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}
	return handleGetResponse{
		Categories: []category{
			{
				Name:          "language",
				SubCategories: []string{"ru-en", "en-ru"},
			},
		},
	}, nil
}
