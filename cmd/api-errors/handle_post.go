package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type handleDataPutRequest struct {
	AppVersion     string         `json:"app_version" validate:"required"`
	OsVersion      string         `json:"os_version" validate:"required"`
	Device         string         `json:"device" validate:"required"`
	ErrorMessage   string         `json:"error_message" validate:"required"`
	ErrorOriginal  string         `json:"error_original" validate:"required"`
	ErrorType      string         `json:"error_type" validate:"required"`
	Timestamp      int            `json:"timestamp" validate:"required"`
	ReplicaID      string         `json:"replica_id" validate:"required"`
	AdditionalInfo map[string]any `json:"additional_info,omitempty"`
}

type handleDataPutResponse struct {
	Status string `json:"status"`
}

func handlePost(ctx context.Context, _ zerolog.Logger, raw json.RawMessage, _ map[string]string) (any, *api.HandleError) {
	var req handleDataPutRequest
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	var (
		key  = time.Now().UTC().Format("logs-2006-01-02.json")
		logs []handleDataPutRequest
	)
	reader, err := s3Bucket.Get(ctx, key, serviceErrorsBucket)
	if err != nil {
		if !errors.Is(err, cloud.ErrBucketObjectNotFound) {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		logs = make([]handleDataPutRequest, 0)
	} else {
		defer reader.Close()
		if err = json.NewDecoder(reader).Decode(&logs); err != nil {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
	}

	logs = append(logs, req)
	data, err := serializer.MarshalJSON(logs)
	if err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	if err = s3Bucket.Put(ctx, key, serviceErrorsBucket, bytes.NewReader(data), cloud.ContentTypeJSON); err != nil {
		return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
	}
	return handleDataPutResponse{Status: "ok"}, nil
}
