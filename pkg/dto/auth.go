package dto

import "github.com/gookit/validate"

// DTOs
type SignInDTO struct {
	Email    string `json:"email" validate:"required|email"`
	Password string `json:"password" validate:"required"`
	Channel  string `json:"channel" validate:"required"`
}

func (dto SignInDTO) Messages() map[string]string {
	return validate.MS{
		"required": "The {field} field is required",
		"email":    "The {field} field must be a valid email address",
	}
}
