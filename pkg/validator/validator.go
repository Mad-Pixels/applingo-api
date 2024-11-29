package validator

import (
	"fmt"
	"strings"

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

func registerCustomTags(v *validator.Validate) {}

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
