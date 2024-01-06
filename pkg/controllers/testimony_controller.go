package controllers

import (
	"strings"

	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TestimonyController struct {
	query *gorm.DB
}

func NewTestimonyController() *TestimonyController {
	q := database.DB.Model(&models.Testimony{}).
		Select("testimonies.*, COUNT(DISTINCT likes.id) AS likes_count, COUNT(DISTINCT comments.id) AS comments_count").
		Joins("JOIN likes ON likes.likeable_id = testimonies.id AND likes.likeable_type = 'testimonies'").
		Joins("JOIN comments ON comments.commentable_id = testimonies.id AND comments.commentable_type = 'testimonies'").
		Group("testimonies.id")

	return &TestimonyController{
		query: q,
	}
}

func (t *TestimonyController) GetAllTestimonies() {}

func (t *TestimonyController) GetTestimonies() {}

func (t *TestimonyController) GetTestimony() {}

// Backoffice routes
func (t *TestimonyController) BackofficeGetTestimonies(c *fiber.Ctx) error {
	var total int64
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)
	search := strings.TrimSpace(c.Query("search", ""))
	query := t.query.Where("testimonies.title LIKE ? OR testimonies.body LIKE ?", "%"+search+"%", "%"+search+"%")

	// Query total
	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Query results
	testimonies := make([]*models.Testimony, 0)
	if err := query.Scopes(models.PaginationScope(c)).Find(&testimonies).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": limit,
		"page":  page,
		"total": total,
		"items": models.TestimoniesToResourceCollection(testimonies),
	})
}

func (t *TestimonyController) BackofficeCreateTestimony(c *fiber.Ctx) error {
	return c.SendString("Ok")
}
