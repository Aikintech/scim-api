package controllers

import (
	"errors"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/aikintech/scim/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	// Parse request
	request := validation.LoginSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Fetch User
	user := models.User{}
	result := config.DB.Where("email = ?", request.Email).First(&user)
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
	accessToken, err := utils.GenerateUserToken(&user, "access", reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	refreshToken, err := utils.GenerateUserToken(&user, "refresh", reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"user": models.AuthUserResource{
				ID:            user.ID,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Email:         user.Email,
				EmailVerified: user.EmailVerifiedAt != nil,
				Avatar:        nil,
				Channels:      user.Channels,
			},
			"tokens": fiber.Map{
				"access":  accessToken,
				"refresh": refreshToken,
			},
		},
	})
}

func Register(c *fiber.Ctx) error {
	// Parse request
	request := validation.RegisterSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Check if email exists
	user := models.User{}
	result := config.DB.Where("email = ?", request.Email).First(&user)
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
	result = config.DB.Model(&models.User{}).Create(&user)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(definitions.MessageResponse{
			Code:    fiber.StatusInternalServerError,
			Message: result.Error.Error(),
		})
	}

	// TODO: Send verification email

	return c.Status(fiber.StatusCreated).JSON(definitions.MessageResponse{
		Code:    fiber.StatusCreated,
		Message: "Account created successfully",
	})
}

func RefreshToken(c *fiber.Ctx) error {
	reference := ulid.Make().String()
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)
	accessToken, err := utils.GenerateUserToken(user, "access", reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}
	//
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": fiber.Map{
			"access": accessToken,
		},
	})
}
