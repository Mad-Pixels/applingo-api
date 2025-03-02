package forge

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
)

// RequestDictionaryCheck is used to perform a dictionary check via an AI model.
// It contains S3 file information and internal buffers to store the downloaded dictionary and processed prompt content.
type RequestDictionaryCheck struct {
	aiModel chatgpt.OpenAIModel

	// DictionaryFile is the S3 key for the dictionary file.
	DictionaryFile string `json:"file"`
	// Prompt is the S3 key (without prefix, just promt name) for the prompt template.
	Prompt *string `json:"prompt"`
	// Model is the user-provided OpenAI model name. It may be invalid.
	Model *string `json:"model"`

	// Internal buffers to hold the fetched data.
	dictionaryBuf *bytes.Buffer
	promptBuf     *bytes.Buffer
}

// NewRequestDictionaryCheck creates a new RequestDictionaryCheck with initialized buffers.
func NewRequestDictionaryCheck() *RequestDictionaryCheck {
	return &RequestDictionaryCheck{
		dictionaryBuf: &bytes.Buffer{},
		promptBuf:     &bytes.Buffer{},
	}
}

// GetPromptBody returns an io.Reader that concatenates the prompt and dictionary content.
// Two newline characters are inserted between the prompt and the dictionary content.
func (r *RequestDictionaryCheck) GetPromptBody() io.Reader {
	return io.MultiReader(r.promptBuf, bytes.NewBufferString("\n\n"), r.dictionaryBuf)
}

// Setup prepares the request and create a prompt for OpenAI or return error.
func (r *RequestDictionaryCheck) Setup(ctx context.Context, s3cli *cloud.Bucket, promptBucketName, dictionaryBucketName string) error {
	var (
		results  = make(chan workerResult, 2)
		setupErr error
		wg       sync.WaitGroup
	)

	if r.Prompt == nil {
		prompt, err := s3cli.GetRandomKey(ctx, promptBucketName, checkPromptPrefix)
		if err != nil {
			return errors.Join(ErrorGetKeyFromBucket(checkPromptPrefix+*r.Prompt, promptBucketName), err)
		}
		r.Prompt = &prompt
	}
	if r.Model == nil {
		model := string(defaultModel)
		r.Model = &model
	} else {
		model, err := chatgpt.ParseModel(aws.ToString(r.Model))
		if err != nil {
			return errors.Join(ErrorOpenAIModelNotSupported(aws.ToString(r.Model)), err)
		}
		r.aiModel = model
	}

	// Operation for getting dictionary content.
	runWorker(ctx, &wg, results, "dictionary", func() error {
		r.dictionaryBuf.Reset()
		if err := s3cli.Read(ctx, r.dictionaryBuf, r.DictionaryFile, dictionaryBucketName); err != nil {
			return errors.Join(ErrorGetBucketFileContent(r.DictionaryFile, dictionaryBucketName), err)
		}
		return nil
	})

	// Operation for getting prompt content.
	runWorker(ctx, &wg, results, "prompt", func() error {
		p := aws.ToString(r.Prompt)

		resp, err := s3cli.GetObjectBody(ctx, p, promptBucketName)
		if err != nil {
			return errors.Join(ErrorGetBucketFileContent(p, promptBucketName), err)
		}
		defer resp.Close()

		r.promptBuf.Reset()
		if err = utils.TemplateFromReaderToWriter(r.promptBuf, resp, r); err != nil {
			return errors.Join(ErrorParseTemplate(p), err)
		}
		return nil
	})

	go func() {
		wg.Wait()
		close(results)
	}()
	for res := range results {
		if res.error != nil {
			workerErr := errors.Join(ErrorWorkerProcess(res.key), res.error)
			if setupErr == nil {
				setupErr = workerErr
			} else {
				setupErr = errors.Join(setupErr, workerErr)
			}
		}
	}
	if setupErr != nil {
		return errors.Join(ErrorSetupProcess, setupErr)
	}
	return nil
}
