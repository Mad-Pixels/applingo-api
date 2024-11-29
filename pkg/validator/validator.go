package validator

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := validator.New()
	registerCustomTags(v)

	return &Validator{validate: v}
}

func (v *Validator) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}

func (v *Validator) ValidateField(field any, tag string) error {
	return v.validate.Var(field, tag)
}

func registerCustomTags(v *validator.Validate) {
	v.RegisterValidation("base_str", func(fl validator.FieldLevel) bool {
		for _, r := range fl.Field().String() {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' && r != '_' && r != ':' {
				return false
			}
		}
		return true
	})

	v.RegisterValidation("ext_str", func(fl validator.FieldLevel) bool {
		validChars := ".-_,#№ +&|[]()\"'{}:"
		for _, r := range fl.Field().String() {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !strings.ContainsRune(validChars, r) {
				return false
			}
		}
		return true
	})

	v.RegisterValidation("lang_code", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if len(s) != 2 {
			return false
		}
		if !unicode.IsLetter(rune(s[0])) || !unicode.IsDigit(rune(s[1])) {
			return false
		}
		return true
	})
}

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
