package goplayground_validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func NewValidate() *validator.Validate {
	validate := validator.New()
	for k, v := range validators() {
		validate.RegisterValidation(k, v.Validate)
	}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	return validate
}
