package trigger

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// HandleFunc is the type for event record handlers.
type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage) error

// Config for api object.
type Config struct {
	MaxWorkers int
}

type trigger struct {
	cfg     Config
	log     zerolog.Logger
	handler HandleFunc
}

// NewLambda creates a new lambda object.
func NewLambda(cfg Config, handler HandleFunc) *trigger {
	return &trigger{
		cfg:     cfg,
		handler: handler,
		log:     logger.InitLogger(),
	}
}

// Handle processes AWS Lambda events, applying the handler function to each record in the event.
// It supports various event types such as DynamoDB and SQS events, and processes records in parallel.
//
// The function unmarshals the raw event, extracts records, and applies the handler to each record
// using a worker pool limited by MaxWorkers. Any errors encountered during processing are collected
// and logged.
//
// Example usage:
//
//	triggerHandler := trigger.NewLambda(
//		trigger.Config{MaxWorkers: 4},
//		func(ctx context.Context, log zerolog.Logger, record json.RawMessage) error {
//			var dynamoRecord events.DynamoDBEventRecord
//			if err := json.Unmarshal(record, &dynamoRecord); err != nil {
//				return errors.Wrap(err, "failed to unmarshal DynamoDB record")
//			}
//			// Process the DynamoDB record
//			return nil
//		},
//	)
//
//	func main() {
//		lambda.Start(triggerHandler.Handle)
//	}
//
// Example event structures:
//
// DynamoDB event:
//
//	{
//		"Records": [
//			{
//				"eventID": "1",
//				"eventName": "INSERT",
//				"eventVersion": "1.0",
//				"eventSource": "aws:dynamodb",
//				"awsRegion": "us-east-1",
//				"dynamodb": {
//					"Keys": {
//						"Id": {
//							"N": "101"
//						}
//					},
//					"NewImage": {
//						"Message": {
//							"S": "New item!"
//						},
//						"Id": {
//							"N": "101"
//						}
//					},
//					"SequenceNumber": "111",
//					"SizeBytes": 26,
//					"StreamViewType": "NEW_AND_OLD_IMAGES"
//				},
//				"eventSourceARN": "arn:aws:dynamodb:us-east-1:account-id:table/ExampleTableWithStream/stream/2015-06-27T00:48:05.899"
//			}
//		]
//	}
//
// SQS event:
//
//	{
//		"Records": [
//			{
//				"messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
//				"receiptHandle": "MessageReceiptHandle",
//				"body": "Hello from SQS!",
//				"attributes": {
//					"ApproximateReceiveCount": "1",
//					"SentTimestamp": "1523232000000",
//					"SenderId": "123456789012",
//					"ApproximateFirstReceiveTimestamp": "1523232000001"
//				},
//				"messageAttributes": {},
//				"md5OfBody": "7b270e59b47ff90a553787216d55d91d",
//				"eventSource": "aws:sqs",
//				"eventSourceARN": "arn:aws:sqs:us-east-1:123456789012:MyQueue",
//				"awsRegion": "us-east-1"
//			}
//		]
//	}
func (t *trigger) Handle(ctx context.Context, event map[string]json.RawMessage) error {
	records, err := t.getRecords(event)
	if err != nil {
		t.log.Error().Err(err).Msg("failed to get records")
		return err
	}
	maxWorkers := t.cfg.MaxWorkers
	if maxWorkers <= 0 {
		maxWorkers = len(records)
	}

	var (
		wg        sync.WaitGroup
		errChan   = make(chan error, len(records))
		semaphore = make(chan struct{}, maxWorkers)
	)

	for _, record := range records {
		wg.Add(1)
		go func(r json.RawMessage) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := t.handler(ctx, t.log, r); err != nil {
				select {
				case errChan <- err:
				default:
					t.log.Error().Err(err).Msg("Error channel full, logging error")
				}
			}
		}(record)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	var hasErrors bool
	for err := range errChan {
		t.log.Error().Err(err).Msg("Error processing record")
		hasErrors = true
	}
	if hasErrors {
		return errors.New("finished with errors")
	}
	return nil
}

func (t *trigger) getRecords(event map[string]json.RawMessage) ([]json.RawMessage, error) {
	if records, ok := event["Records"]; ok {
		var res []json.RawMessage

		if err := serializer.UnmarshalJSON(records, &res); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal records")
		}
		return res, nil
	}
	return nil, errors.New("unsupported event type")
}
