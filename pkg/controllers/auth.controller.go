package controllers

import (
	"github.com/aikintech/scim/pkg/dto"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func SignIn(c *fiber.Ctx) error {
	input := new(dto.SignInDTO)

	// Parse request body
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"message": "An error occurred while parsing your request",
		})
	}

	// Validate request body
	validator := validate.Struct(input)

	if validator.Validate() {
		return c.JSON(input)
	}

	return c.Status(fiber.StatusUnprocessableEntity).JSON(map[string]interface{}{
		"errors": utils.FormatValidationErrors(validator.Errors.All()),
	})
}
