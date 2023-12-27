package controllers

import (
	"fmt"
	"strings"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PodcastController struct {
	query *gorm.DB
}

func NewPodcastController() *PodcastController {
	query := database.DB.Model(&models.Podcast{}).
		Select("podcasts.*, COUNT(DISTINCT likes.id) AS likes_count").
		Joins("LEFT JOIN likes ON likes.likeable_id = podcasts.id AND likes.likeable_type = 'podcasts'").
		Group("podcasts.id")

	return &PodcastController{
		query: query,
	}
}

// ClientListPodcast - List podcasts (paginated)
func (ctrl *PodcastController) ListPodcasts(c *fiber.Ctx) error {
	var total int64
	all := c.Path() == "/podcasts/all"
	sort := c.Query("sort", "newest")
	orderBy := "published_at desc"
	search := strings.Trim(c.Query("search", ""), " ")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if all {
		limit = 999_999_999
	}

	if sort != "newest" {
		orderBy = "published_at asc"
	}

	offset := (page - 1) * limit

	// Fetch podcasts
	podcasts := make([]*models.PodcastResource, 0)
	query := ctrl.query.Where("podcasts.title LIKE ?", "%"+search+"%")

	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := query.
		Offset(offset).
		Limit(limit).
		Order(orderBy).
		Find(&podcasts).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": limit,
		"page":  page,
		"total": total,
		"items": podcasts,
	})
}

// ShowPodcast - Get a podcast
func (ctrl *PodcastController) ShowPodcast(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")

	if len(podcastId) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Message: "No record found",
		})
	}

	// Fetch podcast
	podcast := models.PodcastResource{}
	result := ctrl.query.Where("podcasts.id = ?", podcastId).First(&podcast)

	fmt.Println(podcast)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
				Message: "No record found",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	}

	// Return podcast
	return c.JSON(podcast)
}
