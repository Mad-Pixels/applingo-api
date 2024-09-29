package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/rs/zerolog"
)

type handleDownloadUrlRequest struct {
	DictionaryKey string `json:"dictionary_key" validate:"required,min=4,max=32"`
}

type handleDownloadUrlResponse struct {
	Url string `json:"url"`
}

func handleDownloadUrl(ctx context.Context, _ zerolog.Logger, raw json.RawMessage) (any, *lambda.HandleError) {
	var req handleDownloadUrlRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	url, err := s3Bucket.DownloadUrl(ctx, req.DictionaryKey, serviceDictionaryBucket)
	if err != nil {
		return nil, &lambda.HandleError{Status: http.StatusNotFound, Err: err}
	}
	return handleDownloadUrlResponse{
		Url: url,
	}, nil
}
