package controllers

import (
	"errors"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionController struct{}

func NewTransactionController() *TransactionController {
	return &TransactionController{}
}

// Client handlers
func (t *TransactionController) GetTransactions(c *fiber.Ctx) error {
	var total int64
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

	// Find transactions
	query := database.DB.Model(&models.Transaction{}).Where("user_id = ?", user.ID)
	transactions := make([]*models.Transaction, 0)

	// Total
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := query.Scopes(models.PaginationScope(c)).Preload("User").Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": c.QueryInt("limit", 10),
		"page":  c.QueryInt("page", 1),
		"total": total,
		"items": models.TransactionsToResourceCollection(transactions),
	})
}

func (t *TransactionController) GetTransaction(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	transactionId := c.Params("transactionId", "")

	// Find transaction
	transaction := new(models.Transaction)
	if err := database.DB.Model(&models.Transaction{}).
		Preload("User").
		Where("id = ?", transactionId).
		Where("user_id = ?", user.ID).
		First(&transaction).Error; err != nil {
		status := fiber.StatusNotFound

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusBadRequest
		}
		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(models.TransactionToResource(transaction))
}

func (t *TransactionController) Transact(c *fiber.Ctx) error {
	// Parse request
	request := new(definitions.TransactRequest)
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
