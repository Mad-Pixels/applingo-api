package forge

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
)

// Pkg errors.
var (
	ErrorOpenAIModelNotSupported = func(model string) error {
		return fmt.Errorf("openAI model with name '%s' is not supported, available models: %v", model, chatgpt.AvailableModels())
	}
	ErrorGenerateLanguage = func(direction string) error {
		return fmt.Errorf("failed to generate random language for '%s' direction", direction)
	}
	ErrorGetBucketFileContent = func(key, bucket string) error {
		return fmt.Errorf("failed to read file '%s' from bucket '%s'", key, bucket)
	}
	ErrorGenerateRandomInt = func(min, max int) error {
		return fmt.Errorf("failed to get random int between %d and %d", min, max)
	}
	ErrorGetKeyFromBucket = func(key, bucket string) error {
		return fmt.Errorf("failed to get key '%s' from bucket '%s'", key, bucket)
	}
	ErrorParseTemplate = func(template string) error {
		return fmt.Errorf("failed to parse template '%s'", template)
	}
	ErrorInvalidLanguageCode = func(code string) error {
		return fmt.Errorf("invalid language code '%s'", code)
	}
	ErrorWorkerProcess = func(key string) error {
		return fmt.Errorf("worker '%s' failed", key)
	}

	ErrorGenerateDictionaryDescription = errors.New("failed to generate random dictionary description")
	ErrorSetupProcess                  = errors.New("setup process failed, some workers returned errors")
	ErrorResponseObject                = errors.New("chatgpt response object has incorrect format")
	ErrorGenerateDictionaryTopic       = errors.New("failed to generate random dictionary topic")
	ErrorGenerateLanguageLevel         = errors.New("failed to generate random language level")
	ErrorForgeDictionaryCheck          = errors.New("dictionary check process failed")
	ErrorForgeDictionaryCraft          = errors.New("dictionary craft process failed")
	ErrorReadFromBuffer                = errors.New("failed to read from buffer")
	ErrorWriteToBuffer                 = errors.New("failed write data to buffer")
	ErrorOpenAIProcess                 = errors.New("chatgpt process failed")
)

type workerResult struct {
	key   string
	error error
}

type WorkerFunc func() error

func runWorker(ctx context.Context, wg *sync.WaitGroup, results chan<- workerResult, key string, fn WorkerFunc) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			results <- workerResult{
				key:   key,
				error: ctx.Err(),
			}
			return
		default:
			err := fn()
			results <- workerResult{
				key:   key,
				error: err,
			}
		}
	}()
}
