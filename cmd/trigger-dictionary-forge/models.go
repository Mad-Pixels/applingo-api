package main

// ForgeRequest to lambda.
type ForgeRequest struct {
	OpenAIPromptName      string `json:"openai_prompt_name" validate:"required"`
	OpenAIModelName       string `json:"openai_model_name" validate:"required"`
	DictionaryName        string `json:"dictionary_name" validate:"required"`
	DictionaryTopic       string `json:"dictionary_topic" validate:"required"`
	DictionaryDescription string `json:"dictionary_description" validate:"required"`
	LanguageLevel         string `json:"language_level" validate:"required"`
	LanguageFrom          string `json:"language_from" validate:"required,iso3166_1_alpha2"`
	LanguageTo            string `json:"language_to" validate:"required,iso3166_1_alpha2"`
}

// Message from GPTRequest body.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GPTRequest body.
type GPTRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// GPTResponse body.
type GPTResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
