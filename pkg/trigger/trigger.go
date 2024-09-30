package trigger

import (
	"context"
	"sync"

	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// HandleFunc is the type for event record handlers.
type HandleFunc func(context.Context, zerolog.Logger, any) error

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
// It supports various event types such as DynamoDB and SQS events.
//
// The handler function is called for each record in the event, allowing for parallel processing.
// Any errors encountered during processing are collected and logged.
//
// Example usage:
//
//	triggerHandler := trigger.NewLambda(
//		trigger.Config{
//			MaxWorkers: 10, // Устанавливаем максимальное количество параллельных обработчиков
//		},
//		func(ctx context.Context, log zerolog.Logger, record any) error {
//			switch r := record.(type) {
//			case events.DynamoDBEventRecord:
//				// Process DynamoDB event record
//				return processDynamoDBRecord(r)
//			case events.SQSMessage:
//				// Process SQS message
//				return processSQSMessage(r)
//			default:
//				return errors.New("unknown record type")
//			}
//		},
//	)
//
//	func main() {
//	    lambda.Start(triggerHandler.Handle)
//	}
//
// The Handle function can process various AWS event types. Currently supported:
// - DynamoDB Stream events (*events.DynamoDBEvent)
// - SQS events (*events.SQSEvent)
//
// To add support for more event types, extend the getRecords function accordingly.
func (t *trigger) Handle(ctx context.Context, event any) error {
	records, err := getRecords(event)
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
		errs      = make(chan error, len(records))
		semaphore = make(chan struct{}, maxWorkers)
	)
	for _, recordEvent := range records {
		wg.Add(1)

		go func(e any) {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wg.Done()
			}()

			if err = t.handler(ctx, t.log, e); err != nil {
				select {
				case errs <- err:
				default:
					t.log.Error().Err(err).Msg("Error channel full, logging error")
				}
			}
		}(recordEvent)
	}
	go func() {
		wg.Wait()
		close(errs)
	}()

	var hasErrors bool
	for err = range errs {
		if err != nil {
			t.log.Error().Err(err).Msg("Error processing record")
			hasErrors = true
		}
	}
	if hasErrors {
		return errors.New("finished with errors")
	}
	return nil
}

// getRecords extracts each records from various AWS event types.
// It currently supports DynamoDB and SQS events.
// To add support for more event types, add corresponding case statements.
func getRecords(event any) ([]any, error) {
	switch e := event.(type) {
	case *events.DynamoDBEvent:
		records := make([]any, len(e.Records))
		for i, r := range e.Records {
			records[i] = r
		}
		return records, nil
	case *events.SQSEvent:
		records := make([]any, len(e.Records))
		for i, r := range e.Records {
			records[i] = r
		}
		return records, nil
	default:
		return nil, errors.New("unsupported event type")
	}
}
