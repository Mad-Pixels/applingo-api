package dictionary_check

import (
	"context"
	"sync"

	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/pkg/errors"
)

const (
	defaultModel        = chatgpt.GPT4O
	defaultPromptPrefix = "check"
)

type result struct {
	field string
	value string
	err   error
}

type Result struct {
	Score   int    `json:"score"`
	Message string `json:"message"`
}

func (r *Result) Unmarshal(data []byte) error {
	return serializer.UnmarshalJSON(data, r)
}

func (r *Result) Marshal() ([]byte, error) {
	return serializer.MarshalJSON(r)
}

type Request struct {
	promptBody     string
	dictionaryBody string
	aiModel        chatgpt.OpenAIModel

	File   string `json:"file"`
	Prompt string `json:"prompt"`
	Model  string `json:"model"`
}

// Unmarshal the Request from JSON.
func (r *Request) Unmarshal(data []byte) error {
	return serializer.UnmarshalJSON(data, r)
}

// Marshal the Request to JSON.
func (r *Request) Marshal() ([]byte, error) {
	return serializer.MarshalJSON(r)
}

// GetPromptBody returns the prompt body.
func (r *Request) GetPromptBody() string {
	return r.promptBody + "\n\n" + r.dictionaryBody
}

// GetModel returns the model.
func (r *Request) GetModel() chatgpt.OpenAIModel {
	return r.aiModel
}

// Prepare prepares the request.
func (r *Request) Prepare(ctx context.Context, s3cli *cloud.Bucket, promtBucketName, dictionaryBucketName string) error {
	expectedCount := 2

	var (
		results = make(chan result, expectedCount)
		wg      sync.WaitGroup

		worker = func(field string, f func()) {
			wg.Add(1)
			go func() {
				defer wg.Done()

				select {
				case <-ctx.Done():
					results <- result{field: field, err: ctx.Err()}
					return
				default:
					f()
				}
			}()
		}
	)

	// Operation for getting a random prompt template from S3.
	if r.Prompt == "" {
		prompt, err := s3cli.GetRandomKey(ctx, promtBucketName, defaultPromptPrefix)
		if err != nil {
			return errors.Wrap(err, "failed to get random prompt template from S3")
		}
		r.Prompt = prompt
	}

	// Check model.
	if r.Model == "" {
		r.Model = string(defaultModel)
		r.aiModel = defaultModel
	} else {
		model, err := chatgpt.ParseModel(r.Model)
		if err != nil {
			return errors.Wrap(err, "failed to parse OpenAI model name")
		}
		r.aiModel = model
	}

	// Operation for getting a prompt content.
	worker("prompt", func() {
		tpl, err := s3cli.Read(ctx, r.Prompt, promtBucketName)
		if err != nil {
			results <- result{field: "prompt", err: err}
			return
		}
		prompt, err := utils.Template(string(tpl), r)
		if err != nil {
			results <- result{field: "prompt", value: prompt, err: err}
			return
		}
		r.promptBody = prompt
	})

	// Operation for getting a random dictionary from S3.
	worker("file", func() {
		dictionary, err := s3cli.Read(ctx, r.File, dictionaryBucketName)
		results <- result{field: "file", value: string(dictionary), err: err}
	})

	go func() {
		wg.Wait()
		close(results)
	}()
	for res := range results {
		if res.err != nil {
			return errors.Wrapf(res.err, "failed to prepare %s", res.field)
		}

		switch res.field {
		case "prompt":
			r.promptBody = res.value
		case "file":
			r.dictionaryBody = res.value
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
