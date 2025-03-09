package forge

// RequestDictionaryCheck is used to perform a dictionary check via an AI model.
// It contains information about the S3 file for the dictionary and the prompt template,
// as well as the user-provided model which may need further validation.
type RequestDictionaryCheck struct {
	// PromptName is the name or title of the dictionary check prompt.
	// It serves as an identifier for the particular dictionary check process.
	PromptName *string `json:"prompt_name"`
	// OpenaiModel specifies the OpenAI model to be used for generating the dictionary.
	// This defines the language model that will process the input and produce the output.
	OpenaiModel *string `json:"openai_model"`
	// Temperature controls the creativity or randomness of the language model during generation.
	// A higher temperature value leads to more diverse outputs, while a lower value yields more predictable results.
	Temperature *float64 `json:"temperature"`
}

// NewRequestDictionaryCheck creates a new RequestDictionaryCheck instance.
func NewRequestDictionaryCheck() *RequestDictionaryCheck {
	return &RequestDictionaryCheck{}
}

// Clone creates a deep copy of the RequestDictionaryCheck instance.
// It ensures that all pointer fields are duplicated to prevent accidental modifications to the original instance.
func (r *RequestDictionaryCheck) Clone() *RequestDictionaryCheck {
	clone := NewRequestDictionaryCheck()

	if r.PromptName != nil {
		prompt := *r.PromptName
		clone.PromptName = &prompt
	}
	if r.OpenaiModel != nil {
		model := *r.OpenaiModel
		clone.OpenaiModel = &model
	}
	if r.Temperature != nil {
		temperature := *r.Temperature
		clone.Temperature = &temperature
	}
	return clone
}
