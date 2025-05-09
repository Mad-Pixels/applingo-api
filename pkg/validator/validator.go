// Package validator provides a wrapper around go-playground/validator with custom validation tags.
package validator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Mad-Pixels/applingo-api/pkg/utils"
	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator.Validate instance and exposes custom validation logic.
type Validator struct {
	validate *validator.Validate
}

// New returns a new Validator instance with registered custom validation tags.
func New() *Validator {
	v := validator.New()
	registerCustomTags(v)

	return &Validator{validate: v}
}

// ValidateStruct validates a struct based on its tags.
func (v *Validator) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}

// ValidateField validates a single field against the specified tag.
func (v *Validator) ValidateField(field any, tag string) error {
	return v.validate.Var(field, tag)
}

// StructErrorToString converts validation errors into a human-readable string.
func (v *Validator) StructErrorToString(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var sb strings.Builder
		for _, e := range errs {
			sb.WriteString(fmt.Sprintf("Field '%s' failed validation '%s'\n", e.Field(), e.Tag()))
		}
		return sb.String()
	}
	return err.Error()
}

// registerCustomTags registers project-specific custom validation tags.
func registerCustomTags(v *validator.Validate) {
	_ = v.RegisterValidation("base_str", func(fl validator.FieldLevel) bool {
		invalidChars := "^*%$#@!~`\\/<>?"
		return validateStringWithoutInvalidChars(fl.Field().String(), invalidChars)
	})

	_ = v.RegisterValidation("ext_str", func(fl validator.FieldLevel) bool {
		invalidChars := "^*%$@!~`\\/<>?"
		return validateStringWithoutInvalidChars(fl.Field().String(), invalidChars)
	})

	_ = v.RegisterValidation("lang_code", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if len(s) != 2 {
			return false
		}
		return unicode.IsUpper(rune(s[0])) && unicode.IsDigit(rune(s[1]))
	})

	_ = v.RegisterValidation("file", func(fl validator.FieldLevel) bool {
		return utils.IsFileID(fl.Field().String())
	})
}

// validateStringWithoutInvalidChars checks if the string contains any disallowed characters.
func validateStringWithoutInvalidChars(s string, invalidChars string) bool {
	for _, r := range s {
		if strings.ContainsRune(invalidChars, r) {
			return false
		}
	}
	return true
}
