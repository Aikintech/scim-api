package dto

import (
	"github.com/aikintech/scim/pkg/types"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/gookit/validate"
)

func ValidationMessages() map[string]string {
	return validate.MS{
		"required": "The {field} field is required",
		"email":    "The {field} field must be a valid email address",
		"in":       "The {field} field must be one of the following: {param}",
		"string":   "The {field} field must be of type string",
	}
}

// Validate a data transfer object
// Returns a slice of ValidationErrs if there are any validation errors
func Validate(dto interface{}) types.FormattedValidationErrs {
	validator := validate.Struct(dto)
	var result types.FormattedValidationErrs

	if !validator.Validate() {
		result = utils.FormatValidationErrors(validator.Errors.All())
	}

	return result
}
