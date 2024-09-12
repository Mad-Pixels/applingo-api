package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/rs/zerolog"
)

type handleFilePresignRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=text/csv application/vnd.openxmlformats-officedocument.spreadsheetml.sheet application/vnd.ms-excel"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
}

type handleFilePresignResponse struct {
	Url string `json:"url"`
}

func handleFilePresign(ctx context.Context, _ zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	var req handleFilePresignRequest
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

	url, err := s3Bucket.Presign(ctx, req.Name, serviceProcessingBucket, req.ContentType)
	if err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}
	return handleFilePresignResponse{
		Url: url,
	}, nil
}
