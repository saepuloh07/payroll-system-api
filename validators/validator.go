package validators

import (
	"github.com/go-playground/validator/v10"
)

// Validator global instance
var (
	Validate = validator.New()
)

// ValidateStruct melakukan validasi pada struct
func ValidateStruct(i interface{}) error {
	return Validate.Struct(i)
}
