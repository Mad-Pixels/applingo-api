package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mad-Pixels/applingo-api/openapi-interface"
	"github.com/Mad-Pixels/applingo-api/openapi-interface/gen/applingoapi"
	"github.com/Mad-Pixels/applingo-api/pkg/api"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func handleReportPost(ctx context.Context, _ zerolog.Logger, raw json.RawMessage, _ openapi.QueryParams) (any, *api.HandleError) {
	if !api.MustGetMetaData(ctx).IsDevice() {
		return nil, &api.HandleError{Status: http.StatusForbidden, Err: errors.New("insufficient permissions")}
	}

	var req applingoapi.RequestPostReportV1
	if err := serializer.UnmarshalJSON(raw, &req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}
	if err := validate.ValidateStruct(&req); err != nil {
		return nil, &api.HandleError{Status: http.StatusBadRequest, Err: err}
	}

	var (
		key  = time.Now().UTC().Format("logs-2006-01-02.json")
		logs []applingoapi.RequestPostReportV1
	)
	reader, err := s3Bucket.Get(ctx, key, serviceErrorsBucket)
	if err != nil {
		if !errors.Is(err, cloud.ErrBucketObjectNotFound) {
			return nil, &api.HandleError{Status: http.StatusInternalServerError, Err: err}
		}
		logs = make([]applingoapi.RequestPostReportV1, 0)
	} else {
		defer reader.Close() //nolint:errcheck
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
	return openapi.DataResponseSuccess, nil
}
