package controllers

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/jobs"
	"github.com/aikintech/scim/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ClientListPodcast - List podcasts (paginated)
func ClientListPodcasts(c *fiber.Ctx) error {
	// Validate query params
	sort := c.Query("sort", "newest")
	orderBy := "published_at desc"

	if sort != "newest" {
		orderBy = "published_at asc"
	}

	// Fetch podcasts
	podcasts := make([]models.PodcastResource, 0)
	results := config.DB.Scopes(models.PaginateScope(c)).Model(&models.Podcast{}).Order(orderBy).Find(&podcasts)

	if results.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: results.Error.Error(),
		})
	}

	// Return podcasts
	return c.Status(fiber.StatusOK).JSON(definitions.DataResponse[[]models.PodcastResource]{
		Code: fiber.StatusOK,
		Data: podcasts,
	})
}

// ClientShowPodcast - Get a podcast
func ClientShowPodcast(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")

	if len(podcastId) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Code:    fiber.StatusNotFound,
			Message: "No record found",
		})
	}

	// Fetch podcast
	podcast := models.PodcastResource{}
	result := config.DB.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
				Code:    fiber.StatusNotFound,
				Message: "No record found",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	// Return podcast
	return c.JSON(definitions.DataResponse[models.PodcastResource]{
		Code: fiber.StatusOK,
		Data: podcast,
	})
}

// ClientLikePodcast - Like a podcast
func ClientLikePodcast(c *fiber.Ctx) error {

	return c.SendString("Like podcast")
}

// ClientCommentPodcast - Comment on a podcast
func ClientCommentPodcast(c *fiber.Ctx) error {

	return c.SendString("Like podcast")
}

// ClientUpdatePodcastComment - Update a podcast comment
func ClientUpdatePodcastComment(c *fiber.Ctx) error {

	return c.SendString("Like podcast")
}

// ClientSeedPodcasts - Seed podcasts
func SeedPodcasts(c *fiber.Ctx) error {
	// Background job
	go jobs.SeedPodcasts()

	return c.SendString("Podcasts seeding initiated")
}
