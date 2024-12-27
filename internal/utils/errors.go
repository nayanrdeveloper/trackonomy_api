package utils

import (
	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) map[string]string {
	out := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			fieldName := e.Field()
			out[fieldName] = e.Error()
		}
	} else {
		out["error"] = err.Error()
	}
	return out
}
