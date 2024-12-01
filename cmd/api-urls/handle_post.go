package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/auth"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func handlePost(ctx context.Context, logger zerolog.Logger, body json.RawMessage, _ openapi.QueryParams, reqCtx api.ReqCtx) (any, *api.HandleError) {
	if !reqCtx.HasPermissions(auth.GetPermissionLevel(auth.Device)) {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	var req applingoapi.RequestPostUrlsV1
	if err := serializer.UnmarshalJSON(body, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	switch req.Operation {
	case "upload":
		return handleUpload(ctx, req)
	case "download":
		return handleDownload(ctx, req)
	default:
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: fmt.Errorf("invalid operation")}
	}
}

func handleUpload(ctx context.Context, req applingoapi.RequestPostUrlsV1) (any, *api.HandleError) {
	if req.Identifier == "" {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: errors.New("missing required fields")}
	}
	url, err := s3Bucket.UploadURL(ctx, req.Identifier, serviceProcessingBucket, "text/csv")
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}

	return openapi.DataResponseUrls(applingoapi.UrlsData{
		Url:       url,
		ExpiresIn: 5,
	}), nil
}

func handleDownload(ctx context.Context, req applingoapi.RequestPostUrlsV1) (any, *api.HandleError) {
	if req.Identifier == "" {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: errors.New("missing required fields")}
	}
	url, err := s3Bucket.DownloadURL(ctx, req.Identifier, serviceDictionaryBucket)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusNotFound, Err: err}
	}

	return openapi.DataResponseUrls(applingoapi.UrlsData{
		Url:       url,
		ExpiresIn: 15,
	}), nil
}
