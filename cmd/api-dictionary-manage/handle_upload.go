package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/rs/zerolog"
)

type handleUploadRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=text/csv application/vnd.openxmlformats-officedocument.spreadsheetml.sheet application/vnd.ms-excel"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
}

type handleUploadResponse struct {
	Url string `json:"url"`
}

func handleUpload(ctx context.Context, _ zerolog.Logger, data json.RawMessage) (any, *api.HandleError) {
	var req handleUploadRequest
	if err := serializer.UnmarshalJSON(data, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	url, err := s3Bucket.UploadURL(ctx, req.Name, serviceProcessingBucket, req.ContentType)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleUploadResponse{Url: url}, nil
}
