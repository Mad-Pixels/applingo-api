package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseStrValidation(t *testing.T) {
	validator := New()

	validStrings := []string{
		"simple",
		"with-dash",
		"with_underscore",
		"with.dot",
		"with:colon",
		"mixed-case_with.special:chars123",
		"",
	}

	for _, str := range validStrings {
		err := validator.ValidateField(str, "base_str")
		assert.NoError(t, err, "Expected valid string: %s", str)
	}

	invalidStrings := []string{
		"with space",
		"with*star",
		"with@at",
		"with#hash",
		"with/slash",
		"with\\backslash",
		"with$dollar",
		"with%percent",
		"with&ampersand",
	}

	for _, str := range invalidStrings {
		err := validator.ValidateField(str, "base_str")
		assert.Error(t, err, "Expected invalid string: %s", str)
	}
}

func TestExtStrValidation(t *testing.T) {
	validator := New()

	validStrings := []string{
		"simple",
		// "with.dot",
		// "with:colon",
		// "with space",
		// "with,comma",
		// "with#hash",
		// "with№sign",
		// "with+plus",
		// "with&ampersand",
		// "with|pipe",
		// "with[brackets]",
		// "with(parentheses)",
		// "with\"quotes\"",
		// "with'quotes'",
		// "with{braces}",
	}

	for _, str := range validStrings {
		err := validator.ValidateField(str, "ext_str")
		assert.NoError(t, err, "Expected valid string: %s", str)
	}

	invalidStrings := []string{
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

	for _, str := range invalidStrings {
		err := validator.ValidateField(str, "ext_str")
		assert.Error(t, err, "Expected invalid string: %s", str)
	}
}

func TestLangCodeValidation(t *testing.T) {
	validator := New()

	validCodes := []string{
		"A1", "B2", "C3", "D4", "E5", "Z9",
	}

	for _, code := range validCodes {
		err := validator.ValidateField(code, "lang_code")
		assert.NoError(t, err, "Expected valid code: %s", code)
	}

	invalidCodes := []string{
		"a1",  // lowercase first letter
		"1A",  // digit first
		"AA",  // two letters
		"11",  // two digits
		"A",   // too short
		"A11", // too long
		"",    // empty
	}

	for _, code := range invalidCodes {
		err := validator.ValidateField(code, "lang_code")
		assert.Error(t, err, "Expected invalid code: %s", code)
	}
}

func TestFileValidation(t *testing.T) {
	validator := New()

	validFiles := []string{
		"file.txt",
		"file-name.pdf",
		"file_name.docx",
		"file.name.with.dots",
		"",
	}

	for _, file := range validFiles {
		err := validator.ValidateField(file, "file")
		assert.NoError(t, err, "Expected valid filename: %s", file)
	}

	invalidFiles := []string{
		"file/name.txt",
		"file\\name.txt",
		"file:name.txt",
		"file*name.txt",
		"file name.txt",
		"file,name.txt",
		"file#name.txt",
		"file@name.txt",
	}

	for _, file := range invalidFiles {
		err := validator.ValidateField(file, "file")
		assert.Error(t, err, "Expected invalid filename: %s", file)
	}
}

func TestValidateStringWithChars(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		validChars string
		want       bool
	}{
		{
			name:       "only letters and digits",
			input:      "abc123",
			validChars: ".-_:",
			want:       true,
		},
		{
			name:       "with valid special chars",
			input:      "abc-123.def_ghi:jkl",
			validChars: ".-_:",
			want:       true,
		},
		{
			name:       "with invalid special chars",
			input:      "abc@123",
			validChars: ".-_:",
			want:       false,
		},
		{
			name:       "empty string",
			input:      "",
			validChars: ".-_:",
			want:       true,
		},
		{
			name:       "only special chars",
			input:      ".-_:",
			validChars: ".-_:",
			want:       true,
		},
		{
			name:       "unicode letters",
			input:      "абвгд",
			validChars: ".-_:",
			want:       true,
		},
		{
			name:       "unicode with valid special chars",
			input:      "абвгд-эюя.ёжз",
			validChars: ".-_:",
			want:       true,
		},
		{
			name:       "unicode with invalid special chars",
			input:      "абвгд@эюя",
			validChars: ".-_:",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateStringWithChars(tt.input, tt.validChars)
			assert.Equal(t, tt.want, got)
		})
	}
}
