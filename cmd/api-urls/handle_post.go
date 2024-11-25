package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/v1/urls"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/rs/zerolog"
)

func handlePost(ctx context.Context, logger zerolog.Logger, body json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	var req urls.PostRequest
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	switch req.Operation {
	case urls.OperationUpload:
		return handleUpload(ctx, req)
	case urls.OperationDownload:
		return handleDownload(ctx, req)
	default:
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("invalid operation"),
		}
	}
}

func handleUpload(ctx context.Context, req urls.PostRequest) (any, *api.HandleError) {
	if req.ContentType == "" || req.Name == "" {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("missing required fields"),
		}
	}

	url, err := s3Bucket.UploadURL(ctx, req.Name, serviceProcessingBucket, string(req.ContentType))
	if err != nil {
		return nil, &api.HandleError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return urls.PostResponse{
		URL:       url,
		ExpiresIn: urls.ExpiresIn,
	}, nil
}

func handleDownload(ctx context.Context, req urls.PostRequest) (any, *api.HandleError) {
	if req.Name == "" {
		return nil, &api.HandleError{
			Status: http.StatusBadRequest,
			Err:    fmt.Errorf("missing required fields"),
		}
	}

	url, err := s3Bucket.DownloadURL(ctx, req.Name, serviceDictionaryBucket)
	if err != nil {
		return nil, &api.HandleError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return urls.PostResponse{
		URL:       url,
		ExpiresIn: urls.ExpiresIn,
	}, nil
}
