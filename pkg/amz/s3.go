package amz

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"time"
)

const (
	presignTimeoutMin = 5
)

type S3 struct {
	client *s3.S3
}

// NewS3 ...
func NewS3(session *session.Session) *S3 {
	return &S3{
		client: s3.New(session),
	}
}

// PutRequest ...
func (s *S3) PutRequest(key, bucket, fileType string) (string, error) {
	req, _ := s.client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(fileType),
	})
	return req.Presign(presignTimeoutMin * time.Minute)
}
