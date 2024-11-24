package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/rs/zerolog"
)

type presignedUrlRequest struct {
	ContentType string `json:"content_type,omitempty" validate:"required_if=Operation upload,omitempty,oneof=text/csv application/vnd.openxmlformats-officedocument.spreadsheetml.sheet application/vnd.ms-excel"`
	Name        string `json:"name,omitempty" validate:"required_if=Operation upload,omitempty,min=4,max=32"`
}

type presignedUrlResponse struct {
	Url       string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}

func handlePost(ctx context.Context, _ zerolog.Logger, raw json.RawMessage, params map[string]string) (any, *api.HandleError) {
	// Получаем и проверяем dictionaryId из path параметров
	dictionaryId := params["id"]
	if dictionaryId == "" {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    errors.New("dictionary id is required"),
		}
	}

	// Получаем и проверяем operation из query параметров
	operation := params["operation"]
	if operation == "" {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    errors.New("operation is required"),
		}
	}
	if operation != "upload" && operation != "download" {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    errors.New("operation must be either 'upload' or 'download'"),
		}
	}

	var req presignedUrlRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	switch operation {
	case "download":
		url, err := s3Bucket.DownloadURL(ctx, dictionaryId, serviceDictionaryBucket)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusNotFound, Err: err}
		}
		return presignedUrlResponse{
			Url:       url,
			ExpiresIn: 3600, // 1 hour, настройте в соответствии с вашей конфигурацией
		}, nil

	case "upload":
		url, err := s3Bucket.UploadURL(ctx, req.Name, serviceProcessingBucket, req.ContentType)
		if err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		return presignedUrlResponse{
			Url:       url,
			ExpiresIn: 3600, // 1 hour, настройте в соответствии с вашей конфигурацией
		}, nil

	default:
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    errors.New("invalid operation"),
		}
	}
}
