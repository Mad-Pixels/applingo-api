package forge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
)

const (
	defaultModel = chatgpt.GPT4O

	dictionaryMaxLength = 91
	dictionaryMinLength = 51
	defaultTemperature  = 0.7
	defaultDictionaries = 4
	defaultConcurrent   = 4

	craftPromptPrefix = "craft"
	checkPromptPrefix = "check"
)

// Check sends a request to verify a dictionary using OpenAI and returns a ResponseDictionaryCheck.
// The function performs the following steps:
//  1. Calls Setup on the request to prepare necessary data (e.g., fetching prompt from S3, validating configuration, etc.).
//  2. Reads the prompt content from the request's internal buffer.
//  3. Constructs an OpenAI API request using the prompt and sends it using the chatgpt client.
//  4. Unmarshals the API response into a ResponseDictionaryCheck structure.
//  5. Associates the original request with the response and returns it.
//
// Parameters:
//   - ctx: The context for cancellation and timeouts.
//   - req: A pointer to a RequestDictionaryCheck containing the check parameters.
//   - promptBucket: The name of the S3 bucket where the prompt template is stored.
//   - processingBucket: The name of the S3 bucket for processing data.
//   - chatgptCli: A client for interacting with the OpenAI API.
//   - s3Cli: A client for interacting with the S3 bucket.
//
// Returns:
//   - *ResponseDictionaryCheck: The response containing the check result.
//   - error: An error if any step of the process fails.
func Check(ctx context.Context, req *RequestDictionaryCheck, promptBucket, processingBucket string, chatgptCli *chatgpt.Client, s3Cli *cloud.Bucket) (*ResponseDictionaryCheck, error) {
	if err := req.Setup(ctx, s3Cli, promptBucket, processingBucket); err != nil {
		return nil, errors.Join(ErrorForgeDictionaryCheck, err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		promptStr, err := io.ReadAll(req.GetPromptBody())
		if err != nil {
			return nil, errors.Join(ErrorForgeDictionaryCheck, ErrorReadFromBuffer, err)
		}

		gptReq := chatgpt.NewRequest(
			req.GetModel(),
			[]chatgpt.Message{chatgpt.NewUserMessage(string(promptStr))},
		)
		resp, err := chatgptCli.SendMessage(ctx, gptReq)
		if err != nil {
			return nil, errors.Join(ErrorForgeDictionaryCheck, ErrorOpenAIProcess, err)
		}

		var check ResponseDictionaryCheck
		if err := serializer.UnmarshalJSON([]byte(resp.GetResponseText()), &check); err != nil {
			return nil, errors.Join(ErrorForgeDictionaryCheck, ErrorResponseObject, err)
		}
		check.Request = req
		return &check, nil
	}
}

// Craft sends a request to generate a dictionary using OpenAI and returns a ResponseDictionaryCraft.
// The function performs the following steps:
//  1. Calls Setup on the request to prepare data (fetch prompt from S3, validate model, etc.).
//  2. Reads the prompt content from the request's internal buffer.
//  3. Constructs an OpenAI API request using the prompt and sends it using the chatgpt client.
//  4. Unmarshals the API response into a ResponseDictionaryCraft structure.
//  5. Validates that the generated dictionary contains words and updates the request's word count.
//
// Parameters:
//   - ctx: The context for cancellation and timeouts.
//   - req: A pointer to a RequestDictionaryCraft containing the generation parameters.
//   - promptBucket: The name of the S3 bucket where the prompt template is stored.
//   - chatgptCli: A client for interacting with the OpenAI API.
//   - s3Cli: A client for interacting with the S3 bucket.
//
// Returns:
//   - *ResponseDictionaryCraft: The generated dictionary response.
//   - error: An error if any step of the process fails.
func Craft(ctx context.Context, req *RequestDictionaryCraft, promptBucket string, chatgptCli *chatgpt.Client, s3Cli *cloud.Bucket) (*ResponseDictionaryCraft, error) {
	if err := req.Setup(ctx, s3Cli, promptBucket); err != nil {
		return nil, errors.Join(ErrorForgeDictionaryCraft, err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		promptStr, err := io.ReadAll(req.GetPromptBody())
		if err != nil {
			return nil, errors.Join(ErrorForgeDictionaryCraft, ErrorReadFromBuffer, err)
		}

		gptReq := chatgpt.NewRequest(
			req.GetModel(),
			[]chatgpt.Message{chatgpt.NewUserMessage(string(promptStr))},
		)
		resp, err := chatgptCli.SendMessage(ctx, gptReq)
		if err != nil {
			return nil, errors.Join(ErrorForgeDictionaryCraft, ErrorOpenAIProcess, err)
		}

		var dictionary ResponseDictionaryCraft
		if err := serializer.UnmarshalJSON([]byte(resp.GetResponseText()), &dictionary); err != nil {
			return nil, errors.Join(ErrorForgeDictionaryCraft, ErrorResponseObject, err)
		}
		if len(dictionary.Words) == 0 {
			return nil, errors.Join(ErrorForgeDictionaryCraft, errors.New("dictionary has no words"))
		}

		dictionary.Request = req
		dictionary.Request.WordsCount = len(dictionary.Words)
		return &dictionary, nil
	}
}

// CraftMultiple concurrently generates multiple dictionaries using the Craft function.
// It limits the number of concurrent workers via a semaphore and returns slices of responses and errors.
//
// Parameters:
//   - ctx: The context for cancellation and timeouts.
//   - req: A pointer to the base RequestDictionaryCraft used for cloning each individual request.
//   - promptBucket: The S3 bucket name for the prompt template.
//   - chatgptCli: A client for interacting with the OpenAI API.
//   - s3Cli: A client for interacting with the S3 bucket.
//
// Returns:
//   - []*ResponseDictionaryCraft: A slice of successful dictionary generation responses.
//   - []error: A slice of errors encountered during processing.
func CraftMultiple(ctx context.Context, req *RequestDictionaryCraft, promptBucket string, chatgptCli *chatgpt.Client, s3Cli *cloud.Bucket) ([]*ResponseDictionaryCraft, []error) {
	if req == nil {
		req = NewRequestDictionaryCraft()
	}
	if req.MaxConcurrent <= 0 {
		req.MaxConcurrent = defaultConcurrent
	}
	if req.DictionariesCount < 1 {
		req.DictionariesCount = 1
	}

	var (
		ctxWithCancel, cancel = context.WithCancel(ctx)
		results               = make(chan *ResponseDictionaryCraft, req.DictionariesCount)
		sem                   = make(chan struct{}, req.MaxConcurrent)
		errs                  = make(chan error, req.DictionariesCount)
		wg                    sync.WaitGroup
	)
	defer cancel()

	for i := 0; i < req.DictionariesCount; i++ {
		requestIndex := i
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctxWithCancel.Done():
				errs <- errors.Join(
					ErrorForgeDictionaryCraft,
					fmt.Errorf("request at index %d cancelled by context", requestIndex),
					ctxWithCancel.Err(),
				)
				return
			}

			resp, err := Craft(ctxWithCancel, req.Clone(), promptBucket, chatgptCli, s3Cli)
			if err != nil {
				errs <- errors.Join(
					ErrorForgeDictionaryCraft,
					fmt.Errorf("failed to craft dictionary at index %d", requestIndex),
					err,
				)
				return
			}
			results <- resp
		}()
	}
	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	var (
		dictionaries = make([]*ResponseDictionaryCraft, 0, req.DictionariesCount)
		errorList    = make([]error, 0)
	)
	for resp := range results {
		dictionaries = append(dictionaries, resp)
	}
	for err := range errs {
		errorList = append(errorList, err)
	}
	return dictionaries, errorList
}

// LoadResponseDictionaryCraft loads and unmarshals a ResponseDictionaryCraft from S3.
// It retrieves the object using the provided key and bucket, reads the JSON content,
// and deserializes it into a ResponseDictionaryCraft structure.
func LoadResponseDictionaryCraft(ctx context.Context, s3cli *cloud.Bucket, key, bucket string) (*ResponseDictionaryCraft, error) {
	rc, err := s3cli.GetObjectBody(ctx, key, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get dictionary file %q from bucket %q: %w", key, bucket, err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read dictionary file %q: %w", key, err)
	}

	var dictionary ResponseDictionaryCraft
	if err := serializer.UnmarshalJSON(content, &dictionary); err != nil {
		return nil, fmt.Errorf("failed to unmarshal dictionary file %q: %w", key, err)
	}
	return &dictionary, nil
}
