package utils

import (
	"github.com/aikintech/scim/pkg/types"
	"github.com/matthewhartstonge/argon2"
)

func FormatValidationErrors(errs types.GookitErrs) types.FormattedValidationErrs {
	var formattedErrs types.FormattedValidationErrs

	// Loop over errs
	for key, values := range errs {
		validationError := types.ValidationErr{
			Field:   key,
			Reasons: []string{},
		}

		for _, v := range values {
			validationError.Reasons = append(validationError.Reasons, v)
		}

		formattedErrs = append(formattedErrs, validationError)
	}

	return formattedErrs
}

func HashPassword(password string) string {
	argon := argon2.DefaultConfig()
	hashed, err := argon.HashEncoded([]byte(password))

	if err != nil {
		panic(err)
	}

	return string(hashed)
}

func VerifyPassword(password string, hash string) bool {
	match, err := argon2.VerifyEncoded([]byte(password), []byte(hash))

	if err != nil {
		return false
	}

	return match
}
