package controllers

import (
	"strings"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/jobs"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PodcastController struct{}

func NewPodcastController() *PodcastController {
	return &PodcastController{}
}

// ClientListPodcast - List podcasts (paginated)
func (podCtrl *PodcastController) ListPodcasts(c *fiber.Ctx) error {
	// Validate query params
	sort := c.Query("sort", "newest")
	orderBy := "published_at desc"
	search := strings.Trim(c.Query("search", ""), " ")

	if sort != "newest" {
		orderBy = "published_at asc"
	}

	// Fetch podcasts
	podcasts := make([]models.PodcastResource, 0)
	results := database.DB.Scopes(models.PaginateScope(c)).Model(&models.Podcast{}).Where("title LIKE ?", "%"+search+"%").Order(orderBy).Find(&podcasts)

	if results.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: results.Error.Error(),
		})
	}

	// Return podcasts
	return c.Status(fiber.StatusOK).JSON(podcasts)
}

func (podCtrl *PodcastController) ListAllPodcasts(c *fiber.Ctx) error {
	sort := c.Query("sort", "newest")
	orderBy := "published_at desc"

	if sort != "newest" {
		orderBy = "published_at asc"
	}
	// Fetch podcasts
	podcasts := make([]models.PodcastResource, 0)
	results := database.DB.Model(&models.Podcast{}).Order(orderBy).Find(&podcasts)

	if results.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: results.Error.Error(),
		})
	}

	// Return podcasts
	return c.Status(fiber.StatusOK).JSON(podcasts)
}

// ShowPodcast - Get a podcast
func (podCtrl *PodcastController) ShowPodcast(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")

	if len(podcastId) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Message: "No record found",
		})
	}

	// Fetch podcast
	podcast := models.PodcastResource{}
	result := database.DB.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast)

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

// SeedPodcasts - Seed podcasts
func (podCtrl *PodcastController) SeedPodcasts(c *fiber.Ctx) error {
	// Background job
	go jobs.SeedPodcasts()

	return c.SendString("Podcasts seeding initiated")
}
