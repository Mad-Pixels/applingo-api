package forge

// RequestDictionaryCraft represents a request for generating a crafted dictionary.
// It contains various parameters that define the characteristics of the dictionary to be created.
type RequestDictionaryCraft struct {
	// PromptName is the name or title of the dictionary craft prompt.
	// It serves as an identifier for the particular dictionary creation process.
	PromptName *string `json:"prompt_name"`
	// OpenaiModel specifies the OpenAI model to be used for generating the dictionary.
	// This defines the language model that will process the input and produce the output.
	OpenaiModel *string `json:"openai_model"`
	// DictionaryTopic defines the main topic or subject area for the dictionary.
	// It is used to guide the AI in generating the dictionary.
	DictionaryTopic *string `json:"dictionary_topic"`
	// DictionaryDescription provides a detailed description of the dictionary craft.
	// It is used to guide the AI in generating the dictionary.
	DictionaryDescription *string `json:"dictionary_description"`
	// LanguageLevel indicates the complexity or proficiency level of the language to be used.
	// For example, it could be set to beginner, intermediate, or advanced.
	LanguageLevel *string `json:"language_level"`
	// LanguageFrom denotes the source language from which the dictionary is being crafted.
	// This field is useful in cases where the dictionary involves translation or language comparison.
	LanguageFrom *string `json:"language_from"`
	// LanguageTo represents the target language into which the dictionary content is translated or adapted.
	// It specifies the language that the final dictionary entries will be in.
	LanguageTo *string `json:"language_to"`
	// WordsCount specifies the desired number of words or entries in the generated dictionary.
	// This helps control the length and detail of the dictionary.
	WordsCount *int `json:"words_count"`
	// DictionariesCount determines how many dictionaries should be generated in a single request.
	// This is useful for batch processing or generating multiple variations.
	DictionariesCount *int `json:"dictionaries_count"`
	// MaxConcurrent defines the maximum number of concurrent worker processes to run.
	// It is used to manage the performance and parallel processing of the dictionary generation tasks.
	MaxConcurrent *int `json:"max_concurrent"`
	// Temperature controls the creativity or randomness of the language model during generation.
	// A higher temperature value leads to more diverse outputs, while a lower value yields more predictable results.
	Temperature *float64 `json:"temperature"`
}

// NewDictionaryCraftRequest creates and returns a new instance of RequestDictionaryCraft.
func NewDictionaryCraftRequest() *RequestDictionaryCraft {
	return &RequestDictionaryCraft{}
}

// Clone creates a deep copy of the RequestDictionaryCraft instance.
// It ensures that all pointer fields are duplicated to prevent accidental modifications to the original instance.
func (r *RequestDictionaryCraft) Clone() *RequestDictionaryCraft {
	clone := NewDictionaryCraftRequest()

	if r.PromptName != nil {
		prompt := *r.PromptName
		clone.PromptName = &prompt
	}
	if r.DictionaryTopic != nil {
		topic := *r.DictionaryTopic
		clone.DictionaryTopic = &topic
	}
	if r.DictionaryDescription != nil {
		description := *r.DictionaryDescription
		clone.DictionaryDescription = &description
	}
	if r.LanguageLevel != nil {
		level := *r.LanguageLevel
		clone.LanguageLevel = &level
	}
	if r.LanguageFrom != nil {
		from := *r.LanguageFrom
		clone.LanguageFrom = &from
	}
	if r.LanguageTo != nil {
		to := *r.LanguageTo
		clone.LanguageTo = &to
	}
	if r.Temperature != nil {
		temperature := *r.Temperature
		clone.Temperature = &temperature
	}
	if r.WordsCount != nil {
		wordsCount := *r.WordsCount
		clone.WordsCount = &wordsCount
	}
	if r.DictionariesCount != nil {
		dictionariesCount := *r.DictionariesCount
		clone.DictionariesCount = &dictionariesCount
	}
	if r.MaxConcurrent != nil {
		maxConcurrent := *r.MaxConcurrent
		clone.MaxConcurrent = &maxConcurrent
	}
	if r.OpenaiModel != nil {
		model := *r.OpenaiModel
		clone.OpenaiModel = &model
	}
	return clone
}
