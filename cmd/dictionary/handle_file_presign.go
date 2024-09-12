package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mad-Pixels/lingocards-api/internal/lambda"
	"github.com/Mad-Pixels/lingocards-api/internal/serializer"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"github.com/rs/zerolog"
)

type handleFilePresignRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=text/csv application/vnd.openxmlformats-officedocument.spreadsheetml.sheet application/vnd.ms-excel"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
}

type handleFilePresignResponse struct {
	Url string `json:"url"`
}

func handleFilePresign(_ context.Context, _ zerolog.Logger, data json.RawMessage) (any, *lambda.HandleError) {
	var (
		s3  = cloud.NewBucket(sess)
		req handleFilePresignRequest
	)
	if err := serializer.UnmarshalJSON(data, &req); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}
	if err := validate.Struct(&req); err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	res, err := s3.Presign(req.Name, serviceProcessingBucket, req.ContentType)
	if err != nil {
		return nil, &lambda.HandleError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}
	return handleFilePresignResponse{
		Url: res,
	}, nil
}
