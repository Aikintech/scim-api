package controllers

import (
	"errors"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PrayerController struct{}

func NewPrayerController() *PrayerController {
	return &PrayerController{}
}

func (pryCtrl *PrayerController) MyPrayers(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	sortBy := c.Query("sort", "newest")
	orderBy := "created_at desc"
	if sortBy == "oldest" {
		orderBy = "created_at asc"
	}

	// Get prayer requests
	prayers := []models.PrayerRequestResource{}
	result := database.DB.Model(&models.PrayerRequest{}).Where("user_id = ?", user.ID).Order(orderBy).Find(&prayers)

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

func (pryCtrl *PrayerController) RequestPrayer(c *fiber.Ctx) error {
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
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Create prayer request
	prayer := new(models.PrayerRequest)
	result := database.DB.Model(&prayer).Create(&models.PrayerRequest{
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

	// TODO: Send email to admin

	return c.Status(fiber.StatusCreated).JSON(definitions.MessageResponse{
		Code:    fiber.StatusCreated,
		Message: "Prayer request successfully",
	})
}

// Backoffice handlers
func (pryCtrl *PrayerController) BackOfficeGetPrayers(c *fiber.Ctx) error {
	search := c.Query("search", "")

	// Get prayers
	prayers := make([]*models.PrayerRequest, 0)
	result := database.DB.Scopes(models.PaginateScope(c)).Model(&models.PrayerRequest{}).Preload("User").Where("title LIKE ?", "%"+search+"%").Find(&prayers)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	prayerCollection := models.PrayersToResourceCollection(prayers)

	return c.JSON(definitions.DataResponse[[]models.PrayerRequestResource]{
		Code: fiber.StatusOK,
		Data: prayerCollection,
	})
}