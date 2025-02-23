package cloud

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
)

var (
	ErrBucketObjectNotFound = errors.New("object not found in bucket")
	ErrBucketEmptyKey       = errors.New("empty object key")
	ErrBucketEmptyBucket    = errors.New("empty bucket name")
)

const (
	defaultPartSize   = 5 * 1024 * 1024  // 5MB
	defaultBufferSize = 10 * 1024 * 1024 // 10MB

	uploadTimeout   = 5 * time.Minute
	downloadTimeout = 30 * time.Second
)

const (
	ContentTypeJSON  = "application/json"
	ContentTypeText  = "text/plain"
	ContentTypeHTML  = "text/html"
	ContentTypeCSV   = "text/csv"
	ContentTypePDF   = "application/pdf"
	ContentTypeZIP   = "application/zip"
	ContentTypeImage = "image/jpeg"
)

// Bucket represents an S3 client for object operations.
type Bucket struct {
	client     *s3.Client
	downloader *manager.Downloader
	uploader   *manager.Uploader
}

// NewBucket creates a new instance of S3 client.
func NewBucket(cfg aws.Config) *Bucket {
	client := s3.NewFromConfig(cfg)
	return &Bucket{
		client: client,
		downloader: manager.NewDownloader(client, func(d *manager.Downloader) {
			d.PartSize = defaultPartSize
			d.Concurrency = 10
			d.BufferProvider = manager.NewPooledBufferedWriterReadFromProvider(defaultBufferSize)
		}),
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = defaultPartSize
			u.Concurrency = 10
		}),
	}
}

// validateInput checks basic request parameters.
func validateInput(key, bucket string) error {
	if key == "" {
		return ErrBucketEmptyKey
	}
	if bucket == "" {
		return ErrBucketEmptyBucket
	}
	return nil
}

// UploadURL returns a pre-signed URL for uploading an object to the bucket.
func (b *Bucket) UploadURL(ctx context.Context, key, bucket, contentType string) (string, error) {
	if err := validateInput(key, bucket); err != nil {
		return "", err
	}

	req, err := s3.NewPresignClient(b.client).PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(uploadTimeout))
	if err != nil {
		return "", errors.Wrap(err, "failed to generate upload URL")
	}
	return req.URL, nil
}

// DownloadURL returns a pre-signed URL for downloading an object from the bucket.
func (b *Bucket) DownloadURL(ctx context.Context, key, bucket string) (string, error) {
	if err := validateInput(key, bucket); err != nil {
		return "", err
	}

	req, err := s3.NewPresignClient(b.client).PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(downloadTimeout))
	if err != nil {
		return "", errors.Wrap(err, "failed to generate download URL")
	}
	return req.URL, nil
}

// DownloadToWriter downloads an object from the bucket to an io.Writer.
func (b *Bucket) DownloadToWriter(ctx context.Context, key, bucket string, w io.Writer) error {
	if err := validateInput(key, bucket); err != nil {
		return err
	}
	buffer := manager.NewWriteAtBuffer(make([]byte, 0, defaultBufferSize))

	_, err := b.downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var s3Err *types.NoSuchKey
		if errors.As(err, &s3Err) {
			return ErrBucketObjectNotFound
		}
		return errors.Wrap(err, "failed to download object")
	}
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed write to buffer")
	}
	return nil
}

// Delete removes an object from the bucket.
func (b *Bucket) Delete(ctx context.Context, key, bucket string) error {
	if err := validateInput(key, bucket); err != nil {
		return err
	}

	_, err := b.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var s3Err *types.NoSuchKey
		if errors.As(err, &s3Err) {
			return ErrBucketObjectNotFound
		}
		return errors.Wrap(err, "failed to delete object")
	}
	return nil
}

// Get retrieves an object from the bucket as an io.ReadCloser.
func (b *Bucket) Get(ctx context.Context, key, bucket string) (io.ReadCloser, error) {
	if err := validateInput(key, bucket); err != nil {
		return nil, err
	}

	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var s3Err *types.NoSuchKey
		if errors.As(err, &s3Err) {
			return nil, ErrBucketObjectNotFound
		}
		return nil, errors.Wrap(err, "failed to get object")
	}
	return result.Body, nil
}

// Put uploads an object to the bucket.
func (b *Bucket) Put(ctx context.Context, key, bucket string, body io.Reader, contentType string) error {
	if err := validateInput(key, bucket); err != nil {
		return err
	}

	_, err := b.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload object")
	}
	return nil
}

// Exists checks if an object exists in the bucket.
func (b *Bucket) Exists(ctx context.Context, key, bucket string) (bool, error) {
	if err := validateInput(key, bucket); err != nil {
		return false, err
	}

	_, err := b.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var s3Err *types.NoSuchKey
		if errors.As(err, &s3Err) {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to check object existence")
	}
	return true, nil
}

// Read file from bucket and return content as []byte.
func (b *Bucket) Read(ctx context.Context, key, bucket string) ([]byte, error) {
	if err := validateInput(key, bucket); err != nil {
		return nil, err
	}

	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get object")
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read content from file")
	}
	return data, nil
}

// List returns a list of objects in the bucket.
func (b *Bucket) List(ctx context.Context, bucket string) ([]string, error) {
	if bucket == "" {
		return nil, ErrBucketEmptyBucket
	}

	var keys []string
	paginator := s3.NewListObjectsV2Paginator(b.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list objects")
		}

		for _, obj := range output.Contents {
			keys = append(keys, aws.ToString(obj.Key))
		}
	}

	if len(keys) == 0 {
		return nil, ErrBucketObjectNotFound
	}
	return keys, nil
}
