package controllers

import (
	"errors"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/constants"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	// Parse request
	request := validation.LoginSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Fetch User
	user := models.User{}
	result := database.DB.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid credentials provided",
		})
	}

	// Check password
	ok, err := utils.VerifyPasswordHash(request.Password, user.Password)
	if err != nil || !ok {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid credentials provided",
		})
	}

	// Check if user is verified
	if user.EmailVerifiedAt == nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Account not verified. Please verify your email address",
		})
	}

	// Generate token
	reference := ulid.Make().String()
	accessToken, err := models.GenerateUserToken(user, "access", reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	refreshToken, err := models.GenerateUserToken(user, "refresh", reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": models.UserToResource(&user),
		"tokens": fiber.Map{
			"access":  accessToken,
			"refresh": refreshToken,
		},
	})
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	// Parse request
	request := validation.RegisterSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Check if email exists
	user := models.User{}
	result := database.DB.Where("email = ?", request.Email).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}
	if len(user.ID) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "An account with this email already exists",
		})
	}

	// Create user
	password, err := utils.MakePasswordHash(request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(definitions.MessageResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	user.Email = request.Email
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Password = password
	user.Channels = datatypes.JSON([]byte(`["` + request.Channel + `"]`))
	user.SignUpProvider = "local"
	user.EmailVerifiedAt = nil
	result = database.DB.Model(&models.User{}).Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(definitions.MessageResponse{
			Code:    fiber.StatusInternalServerError,
			Message: result.Error.Error(),
		})
	}

	// TODO: Send verification email

	return c.Status(fiber.StatusCreated).JSON(definitions.SuccessResponse{
		Message: "Account created successfully",
	})
}

func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	reference := ulid.Make().String()
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	accessToken, err := models.GenerateUserToken(*user, "access", reference)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	//
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access": accessToken,
	})
}

func (a *AuthController) ResendEmailVerification(c *fiber.Ctx) error {
	// Parse request
	request := validation.EmailVerificationSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Check if user exists
	message := ""
	user := new(models.User)
	result := database.DB.Model(&models.User{}).Where("email = ?", request.Email).First(&user)

	// Gorm error
	if result.Error != nil {
		message = result.Error.Error()

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = "An account with the selected email does not exist"
		}
	}

	// Sign up provider is not local
	if user.SignUpProvider != "local" {
		message = "Sorry you cannot reset your password with this sign up provider"
	}

	if len(message) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: message,
		})
	}

	// TODO: Send email verification

	return c.JSON(definitions.MessageResponse{
		Code:    fiber.StatusOK,
		Message: "Email verification sent",
	})
}

func (a *AuthController) ForgotPassword(c *fiber.Ctx) error {
	// Parse request
	request := validation.EmailVerificationSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Find user
	user := models.User{}
	result := database.DB.Model(&models.User{}).Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		message := "No account is associated with the email provided"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: message,
		})
	}

	// Check user's sign up provider
	if user.SignUpProvider != "local" {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "The sign up provider for this account does not support password reset",
		})
	}

	// TODO: Send password reset mail
	return c.JSON(definitions.MessageResponse{
		Code:    fiber.StatusOK,
		Message: "Password reset email sent",
	})
}

func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	// Parse request
	request := definitions.ResetPasswordRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Get the user
	user := models.User{}
	result := database.DB.Model(&models.User{}).Where("email ? =", request.Email).First(&user)
	if result.Error != nil {
		message := "No account is associated with the email provided"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: message,
		})
	}

	// Check user's sign up provider
	if user.SignUpProvider != "local" {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "The sign up provider for this account does not support password reset",
		})
	}

	// Set new password
	passwordHash, err := utils.MakePasswordHash(request.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	result = database.DB.Model(&user).Update("Password", passwordHash)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	return c.JSON(definitions.MessageResponse{
		Code:    fiber.StatusOK,
		Message: "Your password has been reset successfully.",
	})
}
