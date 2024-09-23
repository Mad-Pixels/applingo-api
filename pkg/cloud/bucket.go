package cloud

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	presignTimeout = 5
)

type Bucket struct {
	client *s3.Client
}

// NewBucket creates a bucket object.
func NewBucket(cfg aws.Config) *Bucket {
	return &Bucket{
		client: s3.NewFromConfig(cfg),
	}
}

// Presign return url for upload data to bucket.
func (b *Bucket) Presign(ctx context.Context, key, bucket, contentType string) (string, error) {
	presignClient := s3.NewPresignClient(b.client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(presignTimeout*time.Minute))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// Delete an object from the bucket.
func (b *Bucket) Delete(ctx context.Context, key, bucket string) error {
	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return err
}
