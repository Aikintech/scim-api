package controllers

import (
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// Client handlers
func (t *TransactionController) Transact(c *fiber.Ctx) error {
	// Parse request
	request := new(definitions.DonateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}
	return c.SendString("Transact")
}
