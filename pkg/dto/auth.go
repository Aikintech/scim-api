package dto

import (
	"errors"
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/models"
	"gorm.io/gorm"
	"regexp"
)

// LoginDTO
type LoginDTO struct {
	Email    string `json:"email" validate:"required|string|email"`
	Password string `json:"password" validate:"required|string"`
	Channel  string `json:"channel" validate:"required|string|in:web,mobile"`
}

func (dto LoginDTO) Messages() map[string]string {
	return ValidationMessages()
}

// RegisterDTO
type RegisterDTO struct {
	FirstName string `json:"firstName" validate:"required|string|min:3"`
	LastName  string `json:"lastName" validate:"required|string|min:3"`
	Email     string `json:"email" validate:"required|string|email"`
	Password  string `json:"password" validate:"required|string"`
	Channel   string `json:"channel" validate:"required|string|in:web,mobile"`
}

func (dto RegisterDTO) Messages() map[string]string {
	return ValidationMessages()
}

// IsValidPassword Custom validation of checking if password is valid on RegisterDTO
func (dto RegisterDTO) IsValidPassword() bool {
	// At least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	hasLowercase := lowercaseRegex.MatchString(dto.Password)

	// At least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	hasUppercase := uppercaseRegex.MatchString(dto.Password)

	// At least one digit
	digitRegex := regexp.MustCompile(`\d`)
	hasDigit := digitRegex.MatchString(dto.Password)

	// At least one special character (you can customize the characters)
	specialCharRegex := regexp.MustCompile(`[@#$%^&+=!]`)
	hasSpecialChar := specialCharRegex.MatchString(dto.Password)

	// Length between 8 and 30 characters
	lengthRegex := regexp.MustCompile(`^.{8,30}$`)
	hasValidLength := lengthRegex.MatchString(dto.Password)

	// Combine all conditions
	return hasLowercase && hasUppercase && hasDigit && hasSpecialChar && hasValidLength
}

func (dto RegisterDTO) EmailExists() (bool, error) {
	user := new(models.User)

	result := config.DB.Where("email = ?", dto.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Email does not exist
			return false, nil
		}
		// Error occurred
		return false, result.Error
	}
	// Email exists
	return true, nil
}
