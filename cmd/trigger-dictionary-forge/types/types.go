package types

import (
	"context"

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

func (r *Request) Update(
	bucket cloud.Bucket,
	dynamo cloud.Dynamo,
) error {
	if r.Model == "" {
		r.Model = OPENAI_MODEL_DEFAULT
	}

	if r.DictionaryLenght == 0 {
		lenght, err := utils.RandomInt(DICTIONARY_MIN_LENGHT, DICTIONARY_MAX_LENGHT)
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.DictionaryLenght = lenght
	}

	if r.DictionaryName == "" {
		r.DictionaryName = uuid.NewString()
	}

	if r.Prompt == "" {
		prompt, err := bucket.GetRandomKey(context.TODO(), "", "")
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.Prompt = prompt
	}

	if r.DictionaryTopic == "" {
		topic, err := dynamo.GetRandomField(context.TODO(), "", "")
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.DictionaryTopic = topic
	}

	if r.DictionaryDescription == "" {
		description, err := dynamo.GetRandomField(context.TODO(), "", "")
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.DictionaryDescription = description
	}

	if r.LanguageLevel == "" {
		level, err := dynamo.GetRandomField(context.TODO(), "", "")
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.LanguageLevel = level
	}

	if r.LanguageFrom == "" {
		from, err := dynamo.GetRandomField(context.TODO(), "", "")
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.LanguageFrom = from
	}

	if r.LanguageTo == "" {
		to, err := dynamo.GetRandomField(context.TODO(), "", "")
		if err != nil {
			return errors.Wrap(err, "failed to update request")
		}
		r.LanguageTo = to
	}

	return nil
}
