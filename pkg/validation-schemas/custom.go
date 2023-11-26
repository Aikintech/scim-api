package validation_schemas

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func IsValidPasswordValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return isPasswordValid(value)
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
