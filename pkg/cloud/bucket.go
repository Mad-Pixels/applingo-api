package cloud

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	uploadTimeout   = 5
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

// UploadUrl return url for upload data to bucket.
func (b *Bucket) UploadUrl(ctx context.Context, key, bucket, contentType string) (string, error) {
	req, err := s3.NewPresignClient(b.client).PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(uploadTimeout*time.Minute))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// DownloadUrl return url for download data from bucket.
func (b *Bucket) DownloadUrl(ctx context.Context, key, bucket string) (string, error) {
	req, err := s3.NewPresignClient(b.client).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(downloadTimeout*time.Second))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

const (
	defaultPartSize   = 5 * 1024 * 1024  // 5MB
	defaultBufferSize = 10 * 1024 * 1024 // 10MB
)

func (b *Bucket) DownloadToWriter(ctx context.Context, key, bucket string, w io.Writer) error {
	buffer := manager.NewWriteAtBuffer(make([]byte, 0, defaultBufferSize))

	downloader := manager.NewDownloader(b.client, func(d *manager.Downloader) {
		d.PartSize = defaultPartSize
		d.Concurrency = 10
	})

	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	_, err = w.Write(buffer.Bytes())
	return err
}

// Delete an object from the bucket.
func (b *Bucket) Delete(ctx context.Context, key, bucket string) error {
	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var s3Err *types.NoSuchKey
		if !errors.As(err, &s3Err) {
			return err
		}
	}
	return nil
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
