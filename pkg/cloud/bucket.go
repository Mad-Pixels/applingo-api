package cloud

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	presignTimeout  = 5
	downloadTimeout = 30
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

// PresignUrl return url for upload data to bucket.
func (b *Bucket) PresignUrl(ctx context.Context, key, bucket, contentType string) (string, error) {
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

// DownloadUrl return url for download data from bucket.
func (b *Bucket) DownloadUrl(ctx context.Context, key, bucket string) (string, error) {
	presignClient := s3.NewPresignClient(b.client)

	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(downloadTimeout*time.Second))
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

// Get an object from the bucket and returns an io.ReadCloser.
func (b *Bucket) Get(ctx context.Context, key, bucket string) (io.ReadCloser, error) {
	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (b *Bucket) Put(ctx context.Context, key, bucket string, body io.Reader, contentType string) error {
	_, err := b.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}
