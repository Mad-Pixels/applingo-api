package trigger

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Mad-Pixels/lingocards-api/pkg/logger"
	"github.com/Mad-Pixels/lingocards-api/pkg/serializer"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	errUnsupportedEventType = errors.New("unsupported event type")
	errProcessingFailed     = errors.New("finished with errors")
)

const (
	recordsKey = "Records"
)

// HandleFunc is the type for event record handlers.
type HandleFunc func(context.Context, zerolog.Logger, json.RawMessage) error

// Config contains trigger configuration.
type Config struct {
	MaxWorkers int
}

// Trigger handles AWS Lambda events processing.
type Trigger struct {
	cfg     Config
	log     zerolog.Logger
	handler HandleFunc
}

// NewLambda creates a new Lambda trigger instance.
func NewLambda(cfg Config, handler HandleFunc) *Trigger {
	if handler == nil {
		panic("handler function cannot be nil")
	}
	return &Trigger{
		cfg:     cfg,
		handler: handler,
		log:     logger.InitLogger(),
	}
}

// Handle processes AWS Lambda events by applying the handler function to each record.
// It supports various event types such as DynamoDB and SQS events, and processes records in parallel.
func (t *Trigger) Handle(ctx context.Context, event map[string]json.RawMessage) error {
	if ctx == nil {
		ctx = context.Background()
	}
	records, err := t.getRecords(event)
	if err != nil {
		t.log.Error().Err(err).Msg("Failed to get records from event")
		return fmt.Errorf("failed to get records: %w", err)
	}

	if len(records) == 0 {
		t.log.Warn().Msg("No records to process")
		return nil
	}
	maxWorkers := t.getMaxWorkers(len(records))
	return t.processRecords(ctx, records, maxWorkers)
}

// getMaxWorkers determines the number of workers to use.
func (t *Trigger) getMaxWorkers(recordCount int) int {
	if t.cfg.MaxWorkers <= 0 {
		return recordCount
	}
	if t.cfg.MaxWorkers > recordCount {
		return recordCount
	}
	return t.cfg.MaxWorkers
}

// processRecords handles the concurrent processing of records.
func (t *Trigger) processRecords(ctx context.Context, records []json.RawMessage, maxWorkers int) error {
	var (
		wg        sync.WaitGroup
		errChan   = make(chan error, len(records))
		semaphore = make(chan struct{}, maxWorkers)
	)

	t.log.Info().
		Int("total_records", len(records)).
		Int("max_workers", maxWorkers).
		Msg("Starting records processing")

	for i, record := range records {
		wg.Add(1)
		go func(recordNum int, r json.RawMessage) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			recordLogger := t.log.With().Int("record_number", recordNum+1).Logger()

			if err := t.handler(ctx, recordLogger, r); err != nil {
				select {
				case errChan <- fmt.Errorf("record %d processing failed: %w", recordNum+1, err):
				default:
					recordLogger.Error().Err(err).Msg("Error channel full, logging error")
				}
			}
		}(i, record)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	return t.collectErrors(errChan)
}

// collectErrors gathers all errors from the error channel.
func (t *Trigger) collectErrors(errChan <-chan error) error {
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
		t.log.Error().Err(err).Msg("Error processing record")
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %d errors occurred", errProcessingFailed, len(errs))
	}
	return nil
}

// getRecords extracts records from the event payload.
func (t *Trigger) getRecords(event map[string]json.RawMessage) ([]json.RawMessage, error) {
	records, ok := event[recordsKey]
	if !ok {
		return nil, errUnsupportedEventType
	}

	var res []json.RawMessage
	if err := serializer.UnmarshalJSON(records, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal records: %w", err)
	}
	return res, nil
}
