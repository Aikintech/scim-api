package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/definitions"
	v "github.com/go-playground/validator/v10"

	"github.com/iancoleman/strcase"
)

func ValidateStruct(schema interface{}) []definitions.ValidationErr {
	var errs []definitions.ValidationErr
	validate := v.New(v.WithRequiredStructEnabled())

	// Custom error validation registration
	if err := validate.RegisterValidation("validpassword", IsValidPasswordValidation); err != nil {
		fmt.Println("Error registering custom validation :", err.Error())
	}
	if err := validate.RegisterValidation("validfilekey", IsValidUploadFileKey); err != nil {
		fmt.Println("Error registering custom validation :", err.Error())
	}
	if err := validate.RegisterValidation("dateformat", IsValidDateFormat); err != nil {
		fmt.Println("Error registering custom validation :", err.Error())
	}

	// Validate struct
	if err := validate.Struct(schema); err != nil {
		var validationErrors v.ValidationErrors
		errors.As(err, &validationErrors)

		for i, err := range validationErrors {
			field := strcase.ToLowerCamel(err.Field())

			if existsInValidationErrs(field, errs) {
				errs[i].Reasons = append(errs[i].Reasons, getValidationMessage(err))
			} else {
				errs = append(errs, definitions.ValidationErr{Field: field, Reasons: []string{getValidationMessage(err)}})
			}
		}
	}

	return errs
}

func getValidationMessage(err v.FieldError) string {
	// Convert field to camel case
	fieldName := strcase.ToLowerCamel(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", fieldName)
	case "min":
		return fmt.Sprintf("The %s field must be at least %s.", fieldName, err.Param())
	case "max":
		return fmt.Sprintf("The %s field must be at most %s.", fieldName, err.Param())
	case "len":
		return fmt.Sprintf("The %s field must have a length of %s.", fieldName, err.Param())
	case "email":
		return fmt.Sprintf("The %s field must be a valid email.", fieldName)
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL.", fieldName)
	case "numeric":
		return fmt.Sprintf("The %s field must be a numeric value.", fieldName)
	case "alpha":
		return fmt.Sprintf("The %s field must contain only alphabetic characters.", fieldName)
	case "oneof":
		{
			split := strings.Split(err.Param(), " ")
			joined := strings.Join(split, ", ")

			return fmt.Sprintf("The %s field must be one of the following: %s.", fieldName, joined)
		}
	case "validpassword":
		return fmt.Sprintf("The %s field must contain at least one uppercase, one lowercase, one number and one special case character (%s)", fieldName, "@#$%^&+=!")
	case "mimes":
		{
			split := strings.Split(err.Param(), " ")
			joined := strings.Join(split, ", ")

			return fmt.Sprintf("The %s must be one of the following types: %s.", fieldName, joined)
		}
	case "uploadtype":
		{
			stringified := strings.Join(constants.UPLOAD_TYPES, ",")
			return fmt.Sprintf("The %s field must be a valid upload type: %s", fieldName, stringified)
		}
	case "filesize":
		return fmt.Sprintf("The %s field must be a valid file size.", fieldName)
	case "validfilekey":
		return fmt.Sprintf("The %s field provided is invalid.", fieldName)
	case "datetime":
		return fmt.Sprintf("The %s field must be a valid date time.", fieldName)
	case "dateformat":
		{
			params := err.Param()

			return fmt.Sprintf("The %s field must be a valid date time with format %s.", fieldName, params)
		}

	// Add more cases for other validation tags as needed
	default:
		return fmt.Sprintf("Validation failed for %s with tag %s.", fieldName, err.Tag())
	}
}

// TODO: Replace with samber/lo
func existsInValidationErrs(field string, errs []definitions.ValidationErr) bool {
	result := false

	for _, key := range errs {
		if key.Field == field {
			result = true
			break
		}
	}

	return result
}
