package goplayground_validator

import (
	"slices"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func validators() map[string]Validator {
	return map[string]Validator{
		"sex": &sexValidator{},
	}
}

type Validator interface {
	Validate(validator.FieldLevel) bool
	Translate(ut ut.Translator, fe validator.FieldError) string
}

type sexValidator struct{}

func (s *sexValidator) Validate(fl validator.FieldLevel) bool {
	return slices.Contains([]string{"male", "female"}, fl.Field().String())
}

func (s *sexValidator) Translate(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(fe.Tag(), fe.Field())
	return t
}
