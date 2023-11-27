package controllers

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func ClientGetPodcasts(c *fiber.Ctx) error {
	// Validate query params
	sort := c.Query("sort", "newest")
	orderBy := "published_at desc"

	if sort != "newest" {
		orderBy = "published_at asc"
	}

	// Fetch podcasts
	podcasts := make([]models.PodcastResource, 0)
	results := config.DB.Debug().Scopes(models.PaginateScope(c)).Model(&models.Podcast{}).Order(orderBy).Find(&podcasts)

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
