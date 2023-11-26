package controllers

import (
	"errors"
	"fmt"
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/utils"
	validationschemas "github.com/aikintech/scim/pkg/validation-schemas"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	request := new(validationschemas.LoginSchema)

	// Parse request body into LoginSchema
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "",
			Errors:  errs,
		})
	}

	fmt.Println()

	// User
	var user models.User
	result := config.DB.First(&user, "email = ?", request.Email)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Incorrect credentials provided",
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	// Check for hashed password

	return c.JSON(user)
}

func Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
