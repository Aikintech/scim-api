package controllers

import (
	"errors"
	"fmt"
	"strings"

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
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

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

	provider := "stripe"
	if request.Currency == "GHS" && (request.Channel == "mobile_money" || request.Channel == "card") {
		provider = "paystack"
	}

	// Persist transaction
	transaction := models.Transaction{
		UserID:         user.ID,
		Provider:       provider,
		IdempotencyKey: request.IdempotencyKey,
		Currency:       request.Currency,
		Amount:         int64(request.Amount) * 100,
		Type:           request.Type,
		Channel:        request.Channel,
		Description:    strings.TrimSpace(request.Description),
	}

	if err := database.DB.Model(&models.Transaction{}).Create(&transaction).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.TransactionToResource(&transaction))
}

func (t *TransactionController) UpdateTransaction(c *fiber.Ctx) error {
	return nil
}

func (t *TransactionController) PaystackWebhook(c *fiber.Ctx) error {
	// Parse request
	request := definitions.PaystackWebhookPaymentRequest{}
	if err := c.BodyParser(&request); err != nil {
		fmt.Println(err.Error())

		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Println(request)

	return c.SendStatus(fiber.StatusOK)
}

// Backoffice handlers
func (t *TransactionController) BackofficeGetTransactions(c *fiber.Ctx) error {
	var total int64

	// Find transactions
	query := database.DB.Model(&models.Transaction{})
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

func (t *TransactionController) BackofficeGetTransaction(c *fiber.Ctx) error {
	transactionId := c.Params("transactionId", "")

	// Find transaction
	transaction := new(models.Transaction)
	if err := database.DB.Model(&models.Transaction{}).
		Preload("User").
		Where("id = ?", transactionId).
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
