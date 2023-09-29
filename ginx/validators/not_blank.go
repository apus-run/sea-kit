package validators

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

// NotBlank returns whether the given content is not blank
func NotBlank(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		if value != "" && strings.Trim(value, " ") != "" {
			return true
		}
	}

	return false
}
