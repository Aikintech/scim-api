package validation

import (
	"regexp"
	"strings"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/samber/lo"

	"github.com/go-playground/validator/v10"
)

func IsValidPasswordValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return isPasswordValid(value)
}

func IsValidUploadFileKey(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) == 0 {
		return true
	}

	return IsValidFileKey(value)
}

func IsValidDateFormat(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	params := fl.Param()

	if len(params) == 0 {
		return false
	}

	if _, err := time.Parse(params, value); err != nil {
		return false
	}

	return true
}

func isPasswordValid(password string) bool {
	// At least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	hasLowercase := lowercaseRegex.MatchString(password)

	// At least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	hasUppercase := uppercaseRegex.MatchString(password)

	// At least one digit
	digitRegex := regexp.MustCompile(`\d`)
	hasDigit := digitRegex.MatchString(password)

	// At least one special character (you can customize the characters)
	specialCharRegex := regexp.MustCompile(`[@#$%^&+=!]`)
	hasSpecialChar := specialCharRegex.MatchString(password)

	// Length between 8 and 30 characters
	lengthRegex := regexp.MustCompile(`^.{8,30}$`)
	hasValidLength := lengthRegex.MatchString(password)

	// Combine all conditions
	return hasLowercase && hasUppercase && hasDigit && hasSpecialChar && hasValidLength
}

func IsValidFileKey(key string) bool {
	return lo.SomeBy(constants.UPLOAD_TYPES, func(item string) bool {
		uploadType := strings.ToUpper(item)

		return strings.Contains(key, uploadType) // Key contains upload type
	})
}
