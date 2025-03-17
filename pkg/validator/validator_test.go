package validator

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseStrValidation(t *testing.T) {
	var (
		validStrings = []string{
			"simple",
			"with-dash",
			"with_underscore",
			"with.dot",
			"with:colon",
			"mixed-case_with.special:chars123",
			"",
		}
		invalidStrings = []string{
			"with*star",
			"with@at",
			"with/slash",
			"with\\backslash",
			"with$dollar",
			"with%percent",
			"with<arrow",
			"with>arrow",
			"with?question",
		}
	)
	validator := New()

	for _, str := range validStrings {
		err := validator.ValidateField(str, "base_str")
		assert.NoError(t, err, "Expected valid string: %s", str)
	}
	for _, str := range invalidStrings {
		err := validator.ValidateField(str, "base_str")
		assert.Error(t, err, "Expected invalid string: %s", str)
	}
}

func TestExtStrValidation(t *testing.T) {
	var (
		validStrings = []string{
			"simple",
			"with space",
			"with,comma",
			"with#hash",
			"with№sign",
			"with+plus",
			"with&ampersand",
			"with|pipe",
			"with[brackets]",
			"with(parentheses)",
			"with\"quotes\"",
			"with'quotes'",
			"with{braces}",
			"Путеводитель по природе и окружающей среде",
			"日本語テキスト",
			"한국어 텍스트",
			"العربية النص",
			"हिंदी पाठ",
			"Ελληνικά κείμενο",
			"עברית טקסט",
			"Tiếng Việt văn bản",
			"Текст на български",
			"Text with emoji 🚀🌟🌍",
		}

		invalidStrings = []string{
			"with*star",
			"with@at",
			"with/slash",
			"with\\backslash",
			"with%percent",
			"with$dollar",
			"with^caret",
			"with`backtick",
			"with<less than",
			"with>greater than",
		}
	)
	validator := New()

	for _, str := range validStrings {
		err := validator.ValidateField(str, "ext_str")
		assert.NoError(t, err, "Expected valid string: %s", str)
	}
	for _, str := range invalidStrings {
		err := validator.ValidateField(str, "ext_str")
		assert.Error(t, err, "Expected invalid string: %s", str)
	}
}

func TestExtStrWithURLEncodedValues(t *testing.T) {
	var (
		encodedValidStrings = []string{
			"%D0%9F%D1%83%D1%82%D0%B5%D0%B2%D0%BE%D0%B4%D0%B8%D1%82%D0%B5%D0%BB%D1%8C%20%D0%BF%D0%BE%20%D0%BF%D1%80%D0%B8%D1%80%D0%BE%D0%B4%D0%B5%20%D0%B8%20%D0%BE%D0%BA%D1%80%D1%83%D0%B6%D0%B0%D1%8E%D1%89%D0%B5%D0%B9%20%D1%81%D1%80%D0%B5%D0%B4%D0%B5",
			"Dictionary%20with%20spaces",
			"Words%20with%20%23hashtag",
			"Title%20with%20%26ampersand%20symbol",
			"%E6%97%A5%E6%9C%AC%E8%AA%9E%E3%83%86%E3%82%AD%E3%82%B9%E3%83%88",
		}

		encodedInvalidStrings = []string{
			"Invalid%40with%40at",
			"Invalid%25with%25percent",
			"Invalid%24with%24dollar",
		}
	)
	validator := New()

	for _, encodedStr := range encodedValidStrings {
		decodedStr, err := url.QueryUnescape(encodedStr)
		assert.NoError(t, err, "Failed to decode URL: %s", encodedStr)

		err = validator.ValidateField(decodedStr, "ext_str")
		assert.NoError(t, err, "Expected valid decoded string: %s from %s", decodedStr, encodedStr)
	}
	for _, encodedStr := range encodedInvalidStrings {
		decodedStr, err := url.QueryUnescape(encodedStr)
		assert.NoError(t, err, "Failed to decode URL: %s", encodedStr)

		err = validator.ValidateField(decodedStr, "ext_str")
		assert.Error(t, err, "Expected invalid decoded string: %s from %s", decodedStr, encodedStr)
	}
}

func TestLangCodeValidation(t *testing.T) {
	var (
		validCodes = []string{
			"A1", "B2", "C3", "D4", "E5", "Z9",
		}

		invalidCodes = []string{
			"a1",
			"1A",
			"AA",
			"11",
			"A",
			"A11",
			"",
		}
	)
	validator := New()

	for _, code := range validCodes {
		err := validator.ValidateField(code, "lang_code")
		assert.NoError(t, err, "Expected valid code: %s", code)
	}
	for _, code := range invalidCodes {
		err := validator.ValidateField(code, "lang_code")
		assert.Error(t, err, "Expected invalid code: %s", code)
	}
}

func TestFileValidation(t *testing.T) {
	var (
		validFiles = []string{
			"d41d8cd98f00b204e9800998ecf8427e.json",
			"5d41402abc4b2a76b9719d911017c592.json",
			"e10adc3949ba59abbe56e057f20f883e.json",
			"",
		}

		invalidFiles = []string{
			"g41d8cd98f00b204e9800998ecf8427e.json",
			"d41d8cd98f00b204e9800998ecf8427.json",
			"d41d8cd98f00b204e9800998ecf8427eX.json",
			"file.txt",
			"document.pdf",
			"image.png",
			"file-name.pdf",
			"file_name.docx",
			"файл.pdf",
		}
	)
	validator := New()

	for _, file := range validFiles {
		err := validator.ValidateField(file, "file")
		assert.NoError(t, err, "Expected valid filename: %s", file)
	}
	for _, file := range invalidFiles {
		err := validator.ValidateField(file, "file")
		assert.Error(t, err, "Expected invalid filename: %s", file)
	}
}

func TestValidateStringWithoutInvalidChars(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		invalidChars string
		want         bool
	}{
		{
			name:         "only letters and digits",
			input:        "abc123",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         true,
		},
		{
			name:         "with special chars but no invalid ones",
			input:        "abc-123.def_ghi:jkl",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         true,
		},
		{
			name:         "with an invalid char",
			input:        "abc@123",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         false,
		},
		{
			name:         "empty string",
			input:        "",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         true,
		},
		{
			name:         "unicode letters",
			input:        "абвгд",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         true,
		},
		{
			name:         "unicode with spaces",
			input:        "абвгд эюя",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         true,
		},
		{
			name:         "unicode with an invalid char",
			input:        "абвгд@эюя",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         false,
		},
		{
			name:         "Russian text from URL example",
			input:        "Путеводитель по природе и окружающей среде",
			invalidChars: "^*%$@!~`\\/<>?",
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateStringWithoutInvalidChars(tt.input, tt.invalidChars)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRealWorldExamples(t *testing.T) {
	validator := New()

	type paramTest struct {
		name  string
		value string
		tag   string
		valid bool
	}

	tests := []paramTest{
		{
			name:  "Russian dictionary name",
			value: "Путеводитель по природе и окружающей среде",
			tag:   "ext_str",
			valid: true,
		},
		{
			name:  "English author with space",
			value: "Albert Stern",
			tag:   "ext_str",
			valid: true,
		},
		{
			name:  "Language subcategory",
			value: "ru-en",
			tag:   "base_str",
			valid: true,
		},
		{
			name:  "Complex dictionary name with special chars",
			value: "Англо-русский словарь (бизнес & IT)",
			tag:   "ext_str",
			valid: true,
		},
		{
			name:  "Dictionary with quotes",
			value: "Курс \"Разговорный английский\"",
			tag:   "ext_str",
			valid: true,
		},
		{
			name:  "Dictionary with brackets",
			value: "Фразы [начальный уровень]",
			tag:   "ext_str",
			valid: true,
		},
		{
			name:  "Chinese dictionary name",
			value: "汉语词典 (基础)",
			tag:   "ext_str",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateField(tt.value, tt.tag)
			if tt.valid {
				assert.NoError(t, err, "Expected valid string: %s", tt.value)
			} else {
				assert.Error(t, err, "Expected invalid string: %s", tt.value)
			}
		})
	}
}

func TestMockRequestParams(t *testing.T) {
	validator := New()

	type RequestParams struct {
		Name        string `validate:"required,ext_str,min=2,max=36"`
		Author      string `validate:"required,ext_str,min=2,max=24"`
		Subcategory string `validate:"required"`
	}

	t.Run("Actual case with Russian dictionary name", func(t *testing.T) {
		params := RequestParams{
			Name:        "Путеводитель по природе",
			Author:      "Albert Stern",
			Subcategory: "ru-en",
		}
		err := validator.ValidateStruct(params)
		assert.NoError(t, err, "Expected valid params structure")
	})
	t.Run("Case with Chinese dictionary name", func(t *testing.T) {
		params := RequestParams{
			Name:        "自然和环境指南",
			Author:      "李 明",
			Subcategory: "zh-en",
		}
		err := validator.ValidateStruct(params)
		assert.NoError(t, err, "Expected valid params structure with Chinese text")
	})
	t.Run("Case with invalid characters", func(t *testing.T) {
		params := RequestParams{
			Name:        "Guide with@invalid character",
			Author:      "John Doe",
			Subcategory: "en-fr",
		}
		err := validator.ValidateStruct(params)
		assert.Error(t, err, "Expected error for params with invalid characters")
	})
}
