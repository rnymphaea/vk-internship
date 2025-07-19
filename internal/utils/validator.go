package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	cv := &Validator{validator: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return cv
}

func (v *Validator) Validate(s interface{}) error {
	return v.validator.Struct(s)
}

func (v *Validator) FormatValidationErrors(err error) ValidationErrorResponse {
	validationErrors := err.(validator.ValidationErrors)
	errResponse := ValidationErrorResponse{
		Errors: make([]ValidationError, len(validationErrors)),
	}

	for i, fieldError := range validationErrors {
		field := fieldError.Field()
		tag := fieldError.Tag()

		message := generateErrorMessage(field, tag, fieldError.Param())

		errResponse.Errors[i] = ValidationError{
			Field:   field,
			Message: message,
		}
	}

	return errResponse
}

func generateErrorMessage(field, tag, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, param)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
