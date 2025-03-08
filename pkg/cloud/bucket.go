package cloud

import (
	"context"
	"io"
	"math/rand"
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

// Read reads file from bucket and writes content directly to the provided writer.
func (b *Bucket) Read(ctx context.Context, w io.Writer, key, bucket string) error {
	if err := validateInput(key, bucket); err != nil {
		return err
	}

	result, err := b.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return errors.Wrap(err, "failed to get object")
	}
	defer result.Body.Close()

	_, err = io.Copy(w, result.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy content to writer")
	}
	return nil
}

// GetObjectBody returns the response body as io.ReadCloser from S3
func (b *Bucket) GetObjectBody(ctx context.Context, key, bucket string) (io.ReadCloser, error) {
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

	return result.Body, nil
}

// GetRandomKey returns a random key from the bucket.
func (b *Bucket) GetRandomKey(ctx context.Context, bucket, prefix string) (string, error) {
	if bucket == "" {
		return "", ErrBucketEmptyBucket
	}

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int32(1),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	output, err := b.client.ListObjectsV2(ctx, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to get object count")
	}
	if aws.ToInt32(output.KeyCount) == 0 {
		return "", ErrBucketObjectNotFound
	}

	keyCount := aws.ToInt32(output.KeyCount)
	skip := rand.Int31n(keyCount)

	input = &s3.ListObjectsV2Input{
		Bucket:     aws.String(bucket),
		MaxKeys:    aws.Int32(1),
		StartAfter: aws.String(prefix),
	}
	var item *types.Object
	for curr := int32(0); curr <= skip; curr++ {
		output, err := b.client.ListObjectsV2(ctx, input)
		if err != nil {
			return "", errors.Wrap(err, "failed to list objects")
		}
		if len(output.Contents) == 0 {
			break
		}

		item = &output.Contents[0]
		input.StartAfter = item.Key
	}

	if item == nil {
		return "", ErrBucketObjectNotFound
	}
	return aws.ToString(item.Key), nil
}

// Move an object from the source bucket to the destination bucket (or key).
// func (b *Bucket) Move(ctx context.Context, sourceKey, sourceBucket, destKey, destBucket string) error {
// 	if err := validateInput(sourceKey, sourceBucket); err != nil {
// 		return err
// 	}
// 	if err := validateInput(destKey, destBucket); err != nil {
// 		return err
// 	}
// 	copySource := aws.String(sourceBucket + "/" + sourceKey)

// 	_, err := b.client.CopyObject(ctx, &s3.CopyObjectInput{
// 		Bucket:            aws.String(destBucket),
// 		Key:               aws.String(destKey),
// 		CopySource:        copySource,
// 		MetadataDirective: types.MetadataDirectiveCopy,
// 	})
// 	if err != nil {
// 		return errors.Wrap(err, "failed to copy object")
// 	}

// 	waiter := s3.NewObjectExistsWaiter(b.client)
// 	if err = waiter.Wait(ctx, &s3.HeadObjectInput{
// 		Bucket: aws.String(destBucket),
// 		Key:    aws.String(destKey),
// 	}, 30*time.Second); err != nil {
// 		return errors.Wrap(err, "failed to confirm copied object exists")
// 	}
// 	if err = b.Delete(ctx, sourceKey, sourceBucket); err != nil {
// 		return errors.Wrap(err, "failed to delete source object after copy")
// 	}
// 	return nil
// }

func (b *Bucket) Move(ctx context.Context, sourceKey, sourceBucket, destKey, destBucket string) error {
	if err := validateInput(sourceKey, sourceBucket); err != nil {
		return err
	}
	if err := validateInput(destKey, destBucket); err != nil {
		return err
	}
	copySource := aws.String(sourceBucket + "/" + sourceKey)

	// Копируем объект
	_, err := b.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:            aws.String(destBucket),
		Key:               aws.String(destKey),
		CopySource:        copySource,
		MetadataDirective: types.MetadataDirectiveCopy,
	})
	if err != nil {
		return errors.Wrap(err, "failed to copy object")
	}

	// Проверяем наличие объекта с повторными попытками
	retries := 3
	retryDelay := 2 * time.Second

	for i := 0; i < retries; i++ {
		exists, checkErr := b.Exists(ctx, destKey, destBucket)
		if checkErr == nil && exists {
			// Объект существует, можно продолжать
			break
		}

		if i == retries-1 && (checkErr != nil || !exists) {
			// Последняя попытка не удалась
			if checkErr != nil {
				return errors.Wrap(checkErr, "failed to confirm object exists after copy")
			}
			return errors.New("object not found in destination bucket after copy")
		}

		// Ждем перед следующей попыткой
		time.Sleep(retryDelay)
	}

	// Удаляем исходный объект
	if err = b.Delete(ctx, sourceKey, sourceBucket); err != nil {
		return errors.Wrap(err, "failed to delete source object after copy")
	}
	return nil
}
