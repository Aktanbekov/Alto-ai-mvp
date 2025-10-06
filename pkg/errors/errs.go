package errs

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func FromBinding(err error) map[string]string {
	out := map[string]string{"_error": err.Error()}
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		out = map[string]string{}
		for _, fe := range verrs {
			out[fe.Field()] = fe.Tag() // e.g. "required", "email", "min"
		}
	}
	return out
}

// NOTE: Gin pulls in go-playground/validator via binding internally.
