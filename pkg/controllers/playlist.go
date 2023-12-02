package controllers

import (
	"errors"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePlaylist(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)

	// Parse body
	request := new(definitions.StorePlaylistRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Validate body
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&definitions.ValidationErrsResponse{
			Code:    fiber.StatusUnprocessableEntity,
			Message: "Invalid request body",
			Errors:  errs,
		})
	}

	// Create playlist
	playlist := new(models.Playlist)
	podcasts := make([]*models.Podcast, len(request.Podcasts))
	result := config.DB.Find(&podcasts, request.Podcasts)
	if result.Error != nil {
		message := "No podcasts found for the selected ids"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: message,
		})
	}

	playlist.Title = request.Title
	playlist.Description = request.Description
	playlist.UserID = user.ID
	playlist.Podcasts = podcasts

	result = config.DB.Create(&playlist).Save(&playlist)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.PlaylistResource]{
		Code: fiber.StatusCreated,
		Data: *playlist.ToResource(),
	})
}
