package controllers

import (
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/aikintech/scim/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

func RequestPrayer(c *fiber.Ctx) error {
	// Parse body
	var request validation.StorePrayerSchema
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

	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[validation.StorePrayerSchema]{
		Code: fiber.StatusCreated,
		Data: request,
	})
}
