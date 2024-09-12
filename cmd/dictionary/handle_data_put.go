package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/rs/zerolog"
)

type handleDataPutRequest struct {
	Description string `json:"description" validate:"required"`
	Dictionary  string `json:"dictionary" validate:"required"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
	Author      string `json:"author" validate:"required"`
	Category    string `json:"category" validate:"required"`
	SubCategory string `json:"sub_category" validate:"required"`
	Private     bool   `json:"private"`
}

type handleDataPutResponse struct {
	Msg string `json:"msg"`
}

func handleDataPut(_ context.Context, _ zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	var req handleDataPutRequest
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

	return handleDataPutResponse{
		Msg: "OK",
	}, nil
}
