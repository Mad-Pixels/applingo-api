package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/pkg/api"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"

	"github.com/rs/zerolog"
)

type handleDownloadUrlRequest struct {
	Dictionary string `json:"dictionary" validate:"required,min=4,max=32"`
}

type handleDownloadUrlResponse struct {
	Url string `json:"url"`
}

func handleDownloadUrl(ctx context.Context, _ zerolog.Logger, raw json.RawMessage) (any, *api.HandleError) {
	var req handleDownloadUrlRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	url, err := s3Bucket.DownloadURL(ctx, req.Dictionary, serviceDictionaryBucket)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusNotFound, Err: err}
	}
	return handleDownloadUrlResponse{Url: url}, nil
}
