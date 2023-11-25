package controllers

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/dto"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/types"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
)

func Login(c *fiber.Ctx) error {
	// Parse request body
	input := new(dto.LoginDTO)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "An error occurred while parsing your request",
		})
	}

	// Validate request body
	validator := validate.Struct(input)
	if !validator.Validate() {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(types.ValidationErrorResponse{
			Errors: utils.FormatValidationErrors(validator.Errors.All()),
		})
	}

	return c.JSON(input)
}

func Register(c *fiber.Ctx) error {
	input := new(dto.RegisterDTO)

	// Parse request body
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "An error occurred while parsing your request",
		})
	}

	// Validate request body
	if errs := dto.Validate(input); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(types.ValidationErrorResponse{
			Errors: errs,
		})
	}

	// Validate password
	if !input.IsValidPassword() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": append(types.FormattedValidationErrs{}, types.ValidationErr{
				Field:   "password",
				Reasons: []string{"Password must contain at least one uppercase, one lowercase, one number and one special case character"},
			}),
		})
	}

	// Check if email provided exists
	exists, err := input.EmailExists()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(types.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "An account associated with this email exists",
		})
	}

	// Create account
	newUser := models.User{
		ID:             ulid.Make().String(),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       utils.HashPassword(input.Password),
		SignUpProvider: "Local",
		Channels:       datatypes.JSON(`["` + input.Channel + `"]`),
	}

	result := config.DB.Create(&newUser)

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	return c.JSON(types.MessageResponse{
		Code:    fiber.StatusCreated,
		Message: "User account created successfully",
	})
}
