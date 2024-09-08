package main

import (
	"context"
	"encoding/json"
	"github.com/Mad-Pixels/lingocards-api/pkg/cloud"
	"go.uber.org/zap"
)

type handlePresignRequest struct {
	ContentType string `json:"content_type" validate:"required,oneof=text/csv application/vnd.openxmlformats-officedocument.spreadsheetml.sheet application/vnd.ms-excel"`
	Name        string `json:"name" validate:"required,min=4,max=32"`
}

type handlePresignResponse struct {
	Url string `json:"url"`
}

func handlePresign(_ context.Context, _ *zap.Logger, data json.RawMessage) (any, error) {
	var (
		s3  = cloud.NewBucket(sess)
		req handlePresignRequest
	)
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	if err := validate.Struct(&req); err != nil {
		return nil, err
	}

	res, err := s3.Presign(req.Name, serviceProcessingBucket, req.ContentType)
	if err != nil {
		return nil, err
	}
	return handlePresignResponse{
		Url: res,
	}, nil
}
