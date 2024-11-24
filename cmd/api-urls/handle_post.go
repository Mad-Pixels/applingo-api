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

var (
	ErrInvalidDictionaryID = errors.New("invalid dictionary ID")
	ErrInvalidOperation    = errors.New("operation must be either 'upload' or 'download'")
)

// Request models
type GeneratePresignedURLRequest struct {
	Operation   string `json:"operation" validate:"required,oneof=upload download"`
	ContentType string `json:"content_type,omitempty" validate:"required_if=Operation upload,omitempty,oneof=text/csv application/vnd.openxmlformats-officedocument.spreadsheetml.sheet application/vnd.ms-excel"`
	Name        string `json:"name,omitempty" validate:"required_if=Operation upload,omitempty,min=4,max=32"`
}

// Response models
type PresignedURLResponse struct {
	Data struct {
		URL       string `json:"url"`
		ExpiresIn int    `json:"expires_in"`
	} `json:"data"`
}

// Handler represents the dictionary handler
type Handler struct {
	s3Client    S3Client
	logger      zerolog.Logger
	validate    Validator
	downloadTTL int
	uploadTTL   int
}

// S3Client interface defines the required S3 operations
type S3Client interface {
	GenerateDownloadURL(ctx context.Context, key string, bucket string) (string, error)
	GenerateUploadURL(ctx context.Context, key string, bucket string, contentType string) (string, error)
}

// Validator interface for request validation
type Validator interface {
	Struct(interface{}) error
}

// NewHandler creates a new dictionary handler
func NewHandler(s3Client S3Client, logger zerolog.Logger, validator Validator) *Handler {
	return &Handler{
		s3Client:    s3Client,
		logger:      logger,
		validate:    validator,
		downloadTTL: 3600, // 1 hour
		uploadTTL:   3600, // 1 hour
	}
}

// HandleGeneratePresignedURL handles the generation of pre-signed URLs
func (h *Handler) HandleGeneratePresignedURL(ctx context.Context, raw json.RawMessage) (any, *api.HandleError) {
	var req GeneratePresignedURLRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	if err := h.validate.Struct(&req); err != nil {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	var url string
	var err error

	switch req.Operation {
	case "download":
		url, err = h.s3Client.GenerateDownloadURL(ctx, req.Name, "dictionaries")
		if err != nil {
			h.logger.Error().Err(err).Str("dictionary", req.Name).Msg("Failed to generate download URL")
			return nil, &api.HandleError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

	case "upload":
		url, err = h.s3Client.GenerateUploadURL(ctx, req.Name, "uploads", req.ContentType)
		if err != nil {
			h.logger.Error().Err(err).Str("name", req.Name).Msg("Failed to generate upload URL")
			return nil, &api.HandleError{
				Status: http.StatusInternalServerError,
				Err:    err,
			}
		}

	default:
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    ErrInvalidOperation,
		}
	}

	response := PresignedURLResponse{}
	response.Data.URL = url
	response.Data.ExpiresIn = h.downloadTTL

	return response, nil
}
