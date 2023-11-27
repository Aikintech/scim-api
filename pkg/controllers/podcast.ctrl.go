package controllers

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func ClientGetPodcasts(c *fiber.Ctx) error {
	var podcasts []models.PodcastResource

	// Fetch podcasts
	results := config.DB.Debug().Scopes(models.PaginateScope(c)).Model(&models.Podcast{}).Find(&podcasts)

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
