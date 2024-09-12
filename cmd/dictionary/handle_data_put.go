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
}

type handleDataPutResponse struct {
}

func handleDataPut(_ context.Context, _ zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	var (
		req handleDataPutRequest
	)
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

	return handleDataPutResponse{}, nil
}
