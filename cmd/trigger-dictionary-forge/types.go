package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Default OpenAI model.
const (
	OPENAI_MODEL_DEFAULT  = "gpt-3.5-turbo"
	DICTIONARY_MAX_LENGTH = 111
	DICTIONARY_MIN_LENGTH = 51
)

// Request represents the payload for the DictionaryForge function.
// It contains parameters for generating dictionaries, including the OpenAI prompt,
// model, dictionary name, topic, word count, description, language level, and language codes.
type Request struct {
	Prompt                string `json:"prompt"`
	Model                 string `json:"model"`
	DictionaryName        string `json:"dictionary_name"`
	DictionaryTopic       string `json:"dictionary_topic" validate:"min=5,max=500"`
	DictionaryLength      int    `json:"dictionary_length" validate:"min=1,max=500"`
	DictionaryDescription string `json:"dictionary_description" validate:"min=20,max=1000"`
	LanguageLevel         string `json:"language_level"`
	LanguageFrom          string `json:"language_from" validate:"alpha,len=2"`
	LanguageTo            string `json:"language_to" validate:"alpha,len=2"`
}

// Update fills in any missing fields in the Request by generating random values or using default values.
func (r *Request) Update(ctx context.Context, s3cli *cloud.Bucket, bucketName string) error {
	type result struct {
		field string
		value string
		err   error
	}

	var wg sync.WaitGroup
	results := make(chan result, 8)

	// Set default model if not provided.
	if r.Model == "" {
		r.Model = OPENAI_MODEL_DEFAULT
	}

	// Generate a random dictionary length if not provided.
	if r.DictionaryLength == 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			length, err := utils.RandomInt(DICTIONARY_MIN_LENGTH, DICTIONARY_MAX_LENGTH)
			results <- result{field: "length", value: fmt.Sprint(length), err: err}
		}()
	}

	// Generate a UUID for the dictionary name if not provided.
	if r.DictionaryName == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- result{field: "name", value: uuid.NewString(), err: nil}
		}()
	}

	// Retrieve a random prompt from S3 if not provided.
	if r.Prompt == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			prompt, err := s3cli.GetRandomKey(ctx, bucketName, "")
			results <- result{field: "prompt", value: prompt, err: err}
		}()
	}

	// Get a random dictionary topic if not provided.
	if r.DictionaryTopic == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			topic, err := types.GetRandomDictionaryTopic()
			results <- result{field: "topic", value: topic.String(), err: err}
		}()
	}

	// Get a random dictionary description if not provided.
	if r.DictionaryDescription == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			desc, err := types.GetRandomDictionaryDescription()
			results <- result{field: "description", value: desc.String(), err: err}
		}()
	}

	// Get a random language level if not provided.
	if r.LanguageLevel == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			level, err := types.GetRandomLanguageLevel()
			results <- result{field: "level", value: level.String(), err: err}
		}()
	}

	// Select two different random language codes if not provided.
	if r.LanguageFrom == "" || r.LanguageTo == "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			codes := types.AllLanguageCodes()

			idx1, err := utils.RandomInt(0, len(codes)-1)
			if err != nil {
				results <- result{field: "languages", err: err}
				return
			}
			from := codes[idx1]

			remainingCodes := append(codes[:idx1], codes[idx1+1:]...)
			idx2, err := utils.RandomInt(0, len(remainingCodes)-1)
			if err != nil {
				results <- result{field: "languages", err: err}
				return
			}
			to := remainingCodes[idx2]

			results <- result{field: "from", value: from.String(), err: nil}
			results <- result{field: "to", value: to.String(), err: nil}
		}()
	}

	// Close the results channel after all goroutines complete.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Calculate the expected number of results.
	expectedResults := 0
	if r.DictionaryLength == 0 {
		expectedResults++
	}
	if r.DictionaryName == "" {
		expectedResults++
	}
	if r.Prompt == "" {
		expectedResults++
	}
	if r.DictionaryTopic == "" {
		expectedResults++
	}
	if r.DictionaryDescription == "" {
		expectedResults++
	}
	if r.LanguageLevel == "" {
		expectedResults++
	}
	if r.LanguageFrom == "" {
		expectedResults++
	}
	if r.LanguageTo == "" {
		expectedResults++
	}

	// Collect the results.
	receivedResults := 0
	for res := range results {
		if res.err != nil {
			return errors.Wrapf(res.err, "failed to get random %s", res.field)
		}

		switch res.field {
		case "length":
			length, _ := strconv.Atoi(res.value)
			r.DictionaryLength = length
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

		receivedResults++
		if receivedResults == expectedResults {
			break
		}
	}

	// Check if the context has been canceled.
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
