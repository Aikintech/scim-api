package controllers

import (
	"errors"
	"strings"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/ttacon/libphonenumber"
	"gorm.io/gorm"
)

type PrayerController struct{}

func NewPrayerController() *PrayerController {
	return &PrayerController{}
}

func (pryCtrl *PrayerController) MyPrayers(c *fiber.Ctx) error {
	var total int64
	user := c.Locals("user").(*models.User)
	sortBy := c.Query("sort", "newest")
	orderBy := "created_at desc"
	if sortBy == "oldest" {
		orderBy = "created_at asc"
	}

	// Get prayer requests
	query := database.DB.Model(&models.PrayerRequest{}).Where("user_id = ?", user.ID).Count(&total)
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	prayers := make([]*models.PrayerRequest, 0)
	if err := query.Preload("User").Scopes(models.PaginationScope(c)).Order(orderBy).Find(&prayers).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	return c.JSON(definitions.Map{
		"limit": c.QueryInt("limit", 10),
		"page":  c.QueryInt("page", 1),
		"total": total,
		"items": models.PrayersToResourceCollection(prayers),
	})
}

func (pryCtrl *PrayerController) RequestPrayer(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Parse body
	var request definitions.StorePrayerRequest
	if err := c.BodyParser(&request); err != nil {
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

	// Parse phone number
	num, err := libphonenumber.Parse(request.PhoneNumber, strings.ToUpper(request.CountryCode))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: []definitions.ValidationErr{
				{Field: "phoneNumber", Reasons: []string{"Invalid phone number"}},
			},
		})
	}

	// Create prayer request
	prayer := models.PrayerRequest{
		Title:       strings.TrimSpace(request.Title),
		Body:        strings.TrimSpace(request.Description),
		UserID:      user.ID,
		PhoneNumber: libphonenumber.Format(num, libphonenumber.E164),
		CompletedAt: nil,
	}
	result := database.DB.Model(&prayer).Create(&prayer)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// TODO: Send email to admin

	prayer.User = user

	return c.Status(fiber.StatusCreated).JSON(models.PrayerToResource(&prayer))
}

func (pryCtrl *PrayerController) UpdatePrayer(c *fiber.Ctx) error {
	return c.JSON("")
}

// Backoffice handlers
func (pryCtrl *PrayerController) BackOfficeGetPrayers(c *fiber.Ctx) error {
	var total int64
	search := c.Query("search", "")

	// Get total
	query := database.DB.Model(&models.PrayerRequest{}).Where("title LIKE ?", "%"+search+"%")
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Get prayers
	prayers := make([]*models.PrayerRequest, 0)
	result := query.Scopes(models.PaginationScope(c)).Preload("User").Find(&prayers)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	}

	return c.JSON(definitions.Map{
		"limit": c.QueryInt("limit", 10),
		"page":  c.QueryInt("page", 1),
		"total": total,
		"items": models.PrayersToResourceCollection(prayers),
	})
}
