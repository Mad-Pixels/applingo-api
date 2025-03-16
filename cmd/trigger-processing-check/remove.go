package main

import (
	"context"

	"github.com/Mad-Pixels/applingo-api/dynamodb-interface/gen/applingoprocessing"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/aws/aws-lambda-go/events"
)

func remove(ctx context.Context, e events.DynamoDBEventRecord) error {
	var fileID string
	if id, ok := e.Change.Keys[applingoprocessing.ColumnId]; ok {
		fileID = utils.RecordToFileID(id.String())
	}
	if fileID == "" {
		return nil
	}
	return s3Bucket.Delete(ctx, fileID, serviceProcessingBucket)
}
