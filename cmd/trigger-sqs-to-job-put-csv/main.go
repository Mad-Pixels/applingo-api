package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/trigger"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/xuri/excelize/v2"
)

const (
	dictionaryFilenameKey = "filename"
	maxSampleSize         = 1024 * 1024 // 1MB
)

var (
	serviceDictionaryBucket = os.Getenv("SERVICE_DICTIONARY_BUCKET")
	serviceProcessingBucket = os.Getenv("SERVICE_PROCESSING_BUCKET")
	awsRegion               = os.Getenv("AWS_REGION")

	s3Bucket *cloud.Bucket
)

func init() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		panic("unable to load AWS SDK config: " + err.Error())
	}
	s3Bucket = cloud.NewBucket(cfg)
}

func handler(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
	var sqsRecord events.SQSMessage
	if err := serializer.UnmarshalJSON(record, &sqsRecord); err != nil {
		return errors.Wrap(err, "failed to unmarshal SQS record")
	}
	var dynamoDBEvent events.DynamoDBEventRecord
	if err := serializer.UnmarshalJSON([]byte(sqsRecord.Body), &dynamoDBEvent); err != nil {
		return errors.Wrap(err, "failed to unmarshal DynamoDB event from SQS message body")
	}

	bucketKey, ok := dynamoDBEvent.Change.NewImage[dictionaryFilenameKey]
	if !ok {
		return errors.New("'dictionaryFilenameKey' not found in DynamoDB event")
	}
	if bucketKey.DataType() != events.DataTypeString {
		return errors.New("'dictionaryFilenameKey' is not a string in DynamoDB event")
	}
	return processFile(ctx, log, bucketKey.String())
}

func processFile(ctx context.Context, log zerolog.Logger, filename string) error {
	pr, pw := io.Pipe()
	var wg sync.WaitGroup
	var processErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pw.Close()
		if err := s3Bucket.DownloadToWriter(ctx, filename, serviceProcessingBucket, pw); err != nil {
			processErr = errors.Wrapf(err, "failed to download file %s from bucket %s", filename, serviceProcessingBucket)
		}
	}()

	csvData := &strings.Builder{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := convertToCSV(log, pr, csvData); err != nil {
			processErr = errors.Wrap(err, "failed to convert file to CSV")
		}
	}()

	wg.Wait()

	if processErr != nil {
		return processErr
	}

	if err := s3Bucket.Put(ctx, filename, serviceDictionaryBucket, strings.NewReader(csvData.String()), "text/csv"); err != nil {
		return errors.Wrapf(err, "failed to upload file %s to bucket %s", filename, serviceDictionaryBucket)
	}

	return nil
}

func convertToCSV(_ zerolog.Logger, r io.Reader, w io.Writer) error {
	buf := make([]byte, maxSampleSize)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return errors.Wrap(err, "failed to read sample data")
	}
	buf = buf[:n]

	fileType := detectFileType(buf)
	delimiter := detectDelimiter(buf)

	combinedReader := io.MultiReader(bytes.NewReader(buf), r)

	switch fileType {
	case "excel":
		return convertExcelToCSV(combinedReader, w)
	case "csv", "tsv", "custom":
		return convertCSVToCSV(combinedReader, w, delimiter)
	default:
		return errors.New("unsupported file format")
	}
}

func detectFileType(sample []byte) string {
	if len(sample) > 2 && sample[0] == 0x50 && sample[1] == 0x4B {
		return "excel"
	}
	if bytes.Contains(sample, []byte{'\t'}) {
		return "tsv"
	}
	if bytes.Contains(sample, []byte{','}) {
		return "csv"
	}
	return "custom"
}

func detectDelimiter(sample []byte) rune {
	delimiters := []rune{'|', ';', '\t', ','}
	counts := make(map[rune]int)

	for _, d := range delimiters {
		counts[d] = bytes.Count(sample, []byte(string(d)))
	}

	maxCount := 0
	var maxDelimiter rune
	for d, count := range counts {
		if count > maxCount {
			maxCount = count
			maxDelimiter = d
		}
	}

	return maxDelimiter
}

func convertExcelToCSV(r io.Reader, w io.Writer) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return errors.Wrap(err, "failed to open Excel file")
	}
	defer xlsx.Close()

	sheets := xlsx.GetSheetList()
	if len(sheets) == 0 {
		return errors.New("no sheets found in Excel file")
	}

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	rows, err := xlsx.GetRows(sheets[0])
	if err != nil {
		return errors.Wrap(err, "failed to get rows from Excel sheet")
	}

	for _, row := range rows {
		if err := csvWriter.Write(row); err != nil {
			return errors.Wrap(err, "failed to write CSV row")
		}
	}

	return nil
}

func convertCSVToCSV(r io.Reader, w io.Writer, delimiter rune) error {
	reader := csv.NewReader(r)
	reader.Comma = delimiter
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	writer := csv.NewWriter(w)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read CSV record")
		}

		for i, field := range record {
			if strings.Contains(field, ",") && !strings.HasPrefix(field, "\"") {
				record[i] = `"` + strings.Replace(field, `"`, `""`, -1) + `"`
			}
		}

		if err := writer.Write(record); err != nil {
			return errors.Wrap(err, "failed to write CSV record")
		}
	}

	return nil
}

func main() {
	lambda.Start(
		trigger.NewLambda(
			trigger.Config{MaxWorkers: 4},
			handler,
		).Handle,
	)
}
