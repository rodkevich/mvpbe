package validate

import "github.com/go-playground/validator/v10"

// New ...
func New() *validator.Validate {
	return validator.New()
}
