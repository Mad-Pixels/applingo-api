package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/rs/zerolog"
)

type handleDownloadRequest struct {
	Dictionary string `json:"dictionary" validate:"required,min=4,max=40"`
}

type handleDownloadResponse struct {
	Url string `json:"url"`
}

func handleDownload(ctx context.Context, _ zerolog.Logger, raw json.RawMessage) (any, *api.HandleError) {
	var req handleDownloadRequest
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
	return handleDownloadResponse{Url: url}, nil
}
