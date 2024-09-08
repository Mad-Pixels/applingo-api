package cloud

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	presignTimeout = 1
)

type bucket struct {
	client *s3.S3
}

// NewBucket creates a bucket object.
func NewBucket(session *session.Session) *bucket {
	return &bucket{
		client: s3.New(session),
	}
}

// Presign return url for upload data to bucket.
func (b *bucket) Presign(key, bucket, contentType string) (string, error) {
	req, _ := b.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	return req.Presign(presignTimeout * time.Minute)
}
