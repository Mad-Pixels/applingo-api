package forge

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
)

// Package-level errors for various failure scenarios in the dictionary craft and check processes.
var (
	// ErrorOpenAIModelNotSupported indicates that the provided OpenAI model is not supported.
	// It returns an error message including the unsupported model name and a list of available models.
	ErrorOpenAIModelNotSupported = func(model string) error {
		return fmt.Errorf("openAI model with name '%s' is not supported, available models: %v", model, chatgpt.AvailableModels())
	}
	// ErrorGenerateLanguage is returned when a random language generation fails for a given direction (e.g., "from" or "to").
	ErrorGenerateLanguage = func(direction string) error {
		return fmt.Errorf("failed to generate random language for '%s' direction", direction)
	}
	// ErrorGetBucketFileContent indicates failure when attempting to read a file from an S3 bucket.
	ErrorGetBucketFileContent = func(key, bucket string) error {
		return fmt.Errorf("failed to read file '%s' from bucket '%s'", key, bucket)
	}
	// ErrorGenerateRandomInt indicates a failure to generate a random integer within the provided bounds.
	ErrorGenerateRandomInt = func(min, max int) error {
		return fmt.Errorf("failed to get random int between %d and %d", min, max)
	}
	// ErrorGetTemperature indicates that the provided temperature is out of the allowed range.
	ErrorGetTemperature = func(temperature float64, min, max float64) error {
		return fmt.Errorf("temperature must be between %f and %f, got %f", min, max, temperature)
	}
	// ErrorGetWordsCount indicates that the provided word count is not within the allowed range.
	ErrorGetWordsCount = func(words int, min, max int) error {
		return fmt.Errorf("words count must be between %d and %d, got %d", min, max, words)
	}
	// ErrorGetKeyFromBucket indicates a failure to retrieve a specific key from an S3 bucket.
	ErrorGetKeyFromBucket = func(key, bucket string) error {
		return fmt.Errorf("failed to get key '%s' from bucket '%s'", key, bucket)
	}
	// ErrorParseTemplate indicates a failure to parse the specified template.
	ErrorParseTemplate = func(template string) error {
		return fmt.Errorf("failed to parse template '%s'", template)
	}
	// ErrorInvalidLanguageCode is returned when an invalid language code is encountered.
	ErrorInvalidLanguageCode = func(code string) error {
		return fmt.Errorf("invalid language code '%s'", code)
	}
	// ErrorWorkerProcess indicates that a specific worker, identified by its key, has failed.
	ErrorWorkerProcess = func(key string) error {
		return fmt.Errorf("worker '%s' failed", key)
	}

	// Predefined errors for specific process failures.
	ErrorGenerateDictionaryDescription = errors.New("failed to generate random dictionary description")
	ErrorSetupProcess                  = errors.New("setup process failed, some workers returned errors")
	ErrorResponseObject                = errors.New("chatgpt response object has incorrect format")
	ErrorGenerateDictionaryTopic       = errors.New("failed to generate random dictionary topic")
	ErrorGenerateLanguageLevel         = errors.New("failed to generate random language level")
	ErrorForgeDictionaryCheck          = errors.New("dictionary check process failed")
	ErrorForgeDictionaryCraft          = errors.New("dictionary craft process failed")
	ErrorReadFromBuffer                = errors.New("failed to read from buffer")
	ErrorDictionaryFileIsRequired      = errors.New("dictionary file is required")
	ErrorWriteToBuffer                 = errors.New("failed write data to buffer")
	ErrorOpenAIProcess                 = errors.New("chatgpt process failed")
)

// workerResult represents the result from a worker function.
// It contains a key identifying the worker and an error, if any.
type workerResult struct {
	key   string
	error error
}

// WorkerFunc defines the signature for worker functions that return an error.
type WorkerFunc func() error

// runWorker executes a worker function concurrently.
// It adds the worker to the provided WaitGroup, runs the function in a goroutine,
// and sends the result (workerResult) to the provided results channel.
// If the context is canceled, the worker sends the context error instead.
//
// Parameters:
//   - ctx: The context for cancellation and timeouts.
//   - wg: A pointer to a sync.WaitGroup for tracking the worker.
//   - results: A channel to send back the workerResult.
//   - key: A string key identifying this worker.
//   - fn: The worker function to be executed.
func runWorker(ctx context.Context, wg *sync.WaitGroup, results chan<- workerResult, key string, fn WorkerFunc) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		// If the context is canceled, return the context error immediately.
		select {
		case <-ctx.Done():
			results <- workerResult{
				key:   key,
				error: ctx.Err(),
			}
			return
		default:
			// Execute the worker function.
			err := fn()
			results <- workerResult{
				key:   key,
				error: err,
			}
		}
	}()
}
