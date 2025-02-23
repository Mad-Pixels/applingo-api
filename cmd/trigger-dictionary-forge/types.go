package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Mad-Pixels/applingo-api/lingo-interface/types"
	"github.com/Mad-Pixels/applingo-api/pkg/cloud"
	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Default OpenAI model.
const (
	OPENAI_MODEL_DEFAULT  = "gpt-3.5-turbo"
	DICTIONARY_MAX_LENGHT = 111
	DICTIONARY_MIN_LENGHT = 51
)

// Request is the payload for the DictionaryForge lambda function.
// It contains parameters for generating dictionaries, including the OpenAI prompt,
// model, dictionary details, and language settings.
type Request struct {
	// Prompt is the OpenAI prompt name retrieved from the bucket.
	Prompt string `json:"prompt"`
	// Model specifies the OpenAI model to be used for generation.
	Model string `json:"model"`
	// DictionaryName is the name of the dictionary to be crafted.
	DictionaryName string `json:"dictionary_name"`
	// DictionaryTopic indicates the topic of the dictionary.
	DictionaryTopic string `json:"dictionary_topic" validate:"min=5,max=500"`
	// DictionaryLenght represents the count of words in the dictionary.
	DictionaryLenght int `json:"dictionary_lenght" validate:"min=1,max=500"`
	// DictionaryDescription provides a description of the dictionary to be crafted.
	DictionaryDescription string `json:"dictionay_description" validate:"min=20,max=1000"`
	// LanguageLevel denotes the proficiency level of the words in the dictionary.
	LanguageLevel string `json:"language_level"`
	// LanguageFrom is the main language code of the words.
	LanguageFrom string `json:"language_from" validate:"alpha,len=2"`
	// LanguageTo is the language code for the definitions of the words.
	LanguageTo string `json:"language_to" validate:"alpha,len=2"`
}

func (r *Request) Update(ctx context.Context, s3cli *cloud.Bucket, bucketName string) error {
	type result struct {
		field string
		value string
		err   error
	}

	results := make(chan result, 8)

	// Model
	if r.Model == "" {
		r.Model = OPENAI_MODEL_DEFAULT
	}

	// DictionaryLenght
	if r.DictionaryLenght == 0 {
		go func() {
			length, err := utils.RandomInt(DICTIONARY_MIN_LENGHT, DICTIONARY_MAX_LENGHT)
			results <- result{field: "length", value: fmt.Sprint(length), err: err}
		}()
	}

	// DictionaryName
	if r.DictionaryName == "" {
		go func() {
			results <- result{field: "name", value: uuid.NewString(), err: nil}
		}()
	}

	// Prompt
	if r.Prompt == "" {
		go func() {
			prompt, err := s3cli.GetRandomKey(ctx, bucketName, "")
			results <- result{field: "prompt", value: prompt, err: err}
		}()
	}

	// Topic
	if r.DictionaryTopic == "" {
		go func() {
			topic, err := types.GetRandomDictionaryTopic()
			results <- result{field: "topic", value: topic.String(), err: err}
		}()
	}

	// Description
	if r.DictionaryDescription == "" {
		go func() {
			desc, err := types.GetRandomDictionaryDescription()
			results <- result{field: "description", value: desc.String(), err: err}
		}()
	}

	// Level
	if r.LanguageLevel == "" {
		go func() {
			level, err := types.GetRandomLanguageLevel()
			results <- result{field: "level", value: level.String(), err: err}
		}()
	}

	// Languages
	if r.LanguageFrom == "" || r.LanguageTo == "" {
		go func() {
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

			results <- result{
				field: "from",
				value: from.String(),
				err:   nil,
			}
			results <- result{
				field: "to",
				value: to.String(),
				err:   nil,
			}
		}()
	}

	expectedResults := 0
	if r.DictionaryLenght == 0 {
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

	for i := 0; i < expectedResults; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res := <-results:
			if res.err != nil {
				return errors.Wrapf(res.err, "failed to get random %s", res.field)
			}
			switch res.field {
			case "length":
				length, _ := strconv.Atoi(res.value)
				r.DictionaryLenght = length
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
	}
	return nil
}
