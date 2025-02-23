package main

import (
	"context"

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

func (r *Request) Update(ctx context.Context, bucket *cloud.Bucket) error {
	if r.Model == "" {
		r.Model = OPENAI_MODEL_DEFAULT
	}

	if r.DictionaryLenght == 0 {
		lenght, err := utils.RandomInt(DICTIONARY_MIN_LENGHT, DICTIONARY_MAX_LENGHT)
		if err != nil {
			return errors.Wrap(err, "failed to generate random length")
		}
		r.DictionaryLenght = lenght
	}

	if r.DictionaryName == "" {
		r.DictionaryName = uuid.NewString()
	}

	if r.Prompt == "" {
		prompt, err := bucket.GetRandomKey(ctx, "", "")
		if err != nil {
			return errors.Wrap(err, "failed to get random prompt")
		}
		r.Prompt = prompt
	}

	if r.DictionaryTopic == "" {
		topic, err := types.GetRandomDictionaryTopic()
		if err != nil {
			return errors.Wrap(err, "failed to get random topic")
		}
		r.DictionaryTopic = topic.String()
	}

	if r.DictionaryDescription == "" {
		desc, err := types.GetRandomDictionaryDescription()
		if err != nil {
			return errors.Wrap(err, "failed to get random description")
		}
		r.DictionaryDescription = desc.String()
	}

	if r.LanguageLevel == "" {
		level, err := types.GetRandomLanguageLevel()
		if err != nil {
			return errors.Wrap(err, "failed to get random language level")
		}
		r.LanguageLevel = level.String()
	}

	if r.LanguageFrom == "" {
		from, err := types.GetRandomLanguageCode()
		if err != nil {
			return errors.Wrap(err, "failed to get random source language")
		}
		r.LanguageFrom = from.String()
	}

	if r.LanguageTo == "" {
		to, err := types.GetRandomLanguageCode()
		if err != nil {
			return errors.Wrap(err, "failed to get random target language")
		}
		if to.String() == r.LanguageFrom {
			codes := types.AllLanguageCodes()
			for _, code := range codes {
				if code.String() != r.LanguageFrom {
					to = code
					break
				}
			}
		}
		r.LanguageTo = to.String()
	}

	return nil
}
