package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	var errs []string
	for _, v := range e.Errors {
		errs = append(errs, v)
	}
	return strings.Join(errs, ", ")
}

// ParseValidationErrors mengubah validator.ValidationErrors menjadi map field:error
func ParseValidationErrors(err error) ValidationError {
	validationErrors := err.(validator.ValidationErrors)
	errors := make(map[string]string)

	for _, e := range validationErrors {
		tag := e.Tag()
		field := e.Field()

		var message string

		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s", field, e.Param())
		case "gtfield":
			message = fmt.Sprintf("%s must be after %s", field, e.Param())
		default:
			message = fmt.Sprintf("%s is not valid", field)
		}

		errors[field] = message
	}

	return ValidationError{Errors: errors}
}
