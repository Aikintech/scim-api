package controllers

import (
	"errors"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/aikintech/scim/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MyPrayers(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	sortBy := c.Query("sort", "newest")
	orderBy := "created_at desc"
	if sortBy == "oldest" {
		orderBy = "created_at asc"
	}

	// Get prayer requests
	prayers := []models.PrayerRequestResource{}
	result := config.DB.Model(&models.PrayerRequest{}).Where("user_id = ?", user.ID).Order(orderBy).Find(&prayers)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(definitions.MessageResponse{
				Code:    fiber.StatusInternalServerError,
				Message: result.Error.Error(),
			})
		}
	}

	return c.JSON(definitions.DataResponse[[]models.PrayerRequestResource]{
		Code: fiber.StatusOK,
		Data: prayers,
	})
}

func RequestPrayer(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

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

	// Create prayer request
	prayer := new(models.PrayerRequest)
	result := config.DB.Model(&prayer).Create(&models.PrayerRequest{
		Title:       request.Title,
		Body:        request.Description,
		UserID:      user.ID,
		CompletedAt: nil,
	})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(definitions.MessageResponse{
			Code:    fiber.StatusInternalServerError,
			Message: result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(definitions.MessageResponse{
		Code:    fiber.StatusCreated,
		Message: "Prayer request successfully",
	})
}
