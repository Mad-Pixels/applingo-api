package dictionary_craft

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/pkg/chatgpt"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/serializer"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

const (
	DictionaryMaxLength = 91
	DictionaryMinLength = 51
	DefaultModel        = chatgpt.GPT4O
)

// Meta represents metadata about the dictionary.
type Meta struct {
	DictionaryName        string `json:"dictionary_name"`
	DictionaryTopic       string `json:"dictionary_topic"`
	DictionaryDescription string `json:"dictionary_description"`
	LanguageLevel         string `json:"language_level"`
	LanguageFrom          string `json:"language_from"`
	LanguageTo            string `json:"language_to"`
	Author                string `json:"author"`
	WordsCount            int    `json:"words_count"`
}

type result struct {
	field string
	value string
	err   error
}

// Word represents an individual entry in a dictionary.
type Word struct {
	Word        string `json:"word"`
	Translation string `json:"translation"`
	Hint        string `json:"hint"`
	Description string `json:"description"`
}

// Dictionary represents a collection of words with their translations and descriptions.
type Dictionary struct {
	Meta  Meta   `json:"meta"`
	Words []Word `json:"words"`
}

// Unmarshal parses the JSON-encoded data and stores the result in the Dictionary struct.
func (d *Dictionary) Unmarshal(data []byte) error {
	return serializer.UnmarshalJSON(data, d)
}

// Marshal converts the Dictionary into its JSON representation.
func (d *Dictionary) Marshal() ([]byte, error) {
	return serializer.MarshalJSON(d)
}

// Request represents the parameters required to create a dictionary.
type Request struct {
	aiModel    chatgpt.OpenAIModel
	promptBody string

	Prompt                string  `json:"prompt"`
	Model                 string  `json:"model"`
	DictionaryName        string  `json:"dictionary_name"`
	DictionaryTopic       string  `json:"dictionary_topic"`
	DictionaryDescription string  `json:"dictionary_description"`
	LanguageLevel         string  `json:"language_level"`
	LanguageFrom          string  `json:"language_from"`
	LanguageTo            string  `json:"language_to"`
	Temperature           float64 `json:"temperature"`
	DictionaryLength      int     `json:"dictionary_length"`
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
	return r.promptBody
}

// GetModel returns the model.
func (r *Request) GetModel() chatgpt.OpenAIModel {
	return r.aiModel
}

// Prepare fills in any missing fields in the Request with random or default values.
// It concurrently fetches missing values (such as random length, UUID for name,
// prompt template from S3, dictionary topic, description, language level, and language codes)
// and returns an error if any of the operations fail.
func (r *Request) Prepare(ctx context.Context, s3cli *cloud.Bucket, bucketName string) error {
	expectedCount := 0
	if r.DictionaryLength == 0 {
		expectedCount++
	}
	if r.DictionaryName == "" {
		expectedCount++
	}
	if r.Prompt == "" {
		expectedCount++
	}
	if r.DictionaryTopic == "" {
		expectedCount++
	}
	if r.DictionaryDescription == "" {
		expectedCount++
	}
	if r.LanguageLevel == "" {
		expectedCount++
	}
	if r.LanguageFrom == "" || r.LanguageTo == "" {
		expectedCount += 2
	}
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

	// Operation for random dictionary length.
	if r.DictionaryLength == 0 {
		worker("length", func() {
			length, err := utils.RandomInt(DictionaryMinLength, DictionaryMaxLength)
			results <- result{field: "length", value: strconv.Itoa(length), err: err}
		})
	}

	// Operation for generating a UUID for the dictionary name.
	if r.DictionaryName == "" {
		worker("name", func() {
			results <- result{field: "name", value: uuid.NewString(), err: nil}
		})
	}

	// Operation for getting a random template from S3.
	if r.Prompt == "" {
		worker("prompt", func() {
			prompt, err := s3cli.GetRandomKey(ctx, bucketName, "")
			results <- result{field: "prompt", value: prompt, err: err}
		})
	}

	// Operation for getting a random dictionary topic.
	if r.DictionaryTopic == "" {
		worker("topic", func() {
			topic, err := types.GetRandomDictionaryTopic()
			results <- result{field: "topic", value: topic.String(), err: err}
		})
	}

	// Operation for getting a random dictionary description.
	if r.DictionaryDescription == "" {
		worker("description", func() {
			desc, err := types.GetRandomDictionaryDescription()
			results <- result{field: "description", value: desc.String(), err: err}
		})
	}

	// Operation for getting a random language level.
	if r.LanguageLevel == "" {
		worker("level", func() {
			level, err := types.GetRandomLanguageLevel()
			results <- result{field: "level", value: level.String(), err: err}
		})
	}

	// Operation for selecting two different random language codes.
	if r.LanguageFrom == "" || r.LanguageTo == "" {
		worker("languages", func() {
			codes := types.AllLanguageCodes()
			idx1, err := utils.RandomInt(0, len(codes)-1)
			if err != nil {
				results <- result{field: "languages", err: err}
			}

			remainingCodes := slices.Delete(codes, idx1, idx1+1)
			idx2, err := utils.RandomInt(0, len(remainingCodes)-1)
			if err != nil {
				results <- result{field: "languages", err: err}
			}

			results <- result{field: "from", value: codes[idx1].String(), err: nil}
			results <- result{field: "to", value: remainingCodes[idx2].String(), err: nil}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()
	for res := range results {
		if res.err != nil {
			return errors.Wrapf(res.err, "failed to get random %s", res.field)
		}

		switch res.field {
		case "length":
			val, convErr := strconv.Atoi(res.value)
			if convErr != nil {
				return errors.Wrap(convErr, "failed to convert length")
			}
			r.DictionaryLength = val
		case "name":
			r.DictionaryName = res.value
		case "prompt":
			r.Prompt = res.value
		case "topic":
			r.DictionaryTopic = res.value
		case "description":
			r.DictionaryDescription = res.value
		case "level":
			r.LanguageLevel = res.value
		case "from":
			r.LanguageFrom = res.value
		case "to":
			r.LanguageTo = res.value
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Prepare prompt content.
		tpl, err := s3cli.Read(ctx, r.Prompt, bucketName)
		if err != nil {
			return errors.Wrap(err, "failed to get prompt template data")
		}
		prompt, err := utils.Template(string(tpl), r)
		if err != nil {
			return errors.Wrap(err, "failed to prepare prompt")
		}
		r.promptBody = prompt

		// Check Model.
		if r.Model == "" {
			r.Model = string(DefaultModel)
			r.aiModel = DefaultModel
		} else {
			model, err := chatgpt.ParseModel(r.Model)
			if err != nil {
				return fmt.Errorf("model %s not found in available models list: %v", r.Model, chatgpt.AvailableModels())
			}
			r.aiModel = model
		}
		return nil
	}
}
