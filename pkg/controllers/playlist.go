package controllers

import (
	"errors"
	"fmt"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PlaylistController struct{}

func NewPlaylistController() *PlaylistController {
	return &PlaylistController{}
}

// GetPlaylists - Get user playlists
func (plCtrl *PlaylistController) GetPlaylists(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)

	// Get playlists
	playlists := make([]*models.Playlist, 0)
	if result := config.DB.Preload("Podcasts").Where("user_id = ?", user.ID).Find(&playlists); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	// Convert to resource
	resources := make([]*models.PlaylistResource, len(playlists))
	for i, playlist := range playlists {
		resources[i] = playlist.ToResource()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": resources,
	})
}

// CreatePlaylist - Create a playlist
func (plCtrl *PlaylistController) CreatePlaylist(c *fiber.Ctx) error {
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
	podcasts := make([]models.Podcast, 0)
	if result := config.DB.Find(&podcasts, request.Podcasts); result.Error != nil {
		message := "No podcasts found for the selected ids"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: message,
		})
	}

	playlist := models.Playlist{
		Title:       request.Title,
		Description: request.Description,
		UserID:      user.ID,
		Podcasts:    podcasts,
	}

	if result := config.DB.Create(&playlist); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.PlaylistResource]{
		Code: fiber.StatusCreated,
		Data: *playlist.ToResource(),
	})
}

// GetPlaylist - Get a playlist
func (plCtrl *PlaylistController) GetPlaylist(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)

	// Get playlist
	playlist := new(models.Playlist)
	if result := config.DB.Preload("Podcasts").Where("id = ? AND user_id = ?", c.Params("playlistId"), user.ID).First(&playlist); result.Error != nil {
		message := "Playlist not found"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Code:    fiber.StatusNotFound,
			Message: message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(definitions.DataResponse[models.PlaylistResource]{
		Code: fiber.StatusOK,
		Data: *playlist.ToResource(),
	})
}

// TODO: Fix me
func (plCtrl *PlaylistController) UpdatePlaylist(c *fiber.Ctx) error {
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

	// Get playlist
	trx := config.DB.Begin()
	playlist := models.Playlist{}
	result := trx.Preload("Podcasts").Where("id = ? AND user_id = ?", c.Params("playlistId"), user.ID).First(&playlist)
	if result.Error != nil {
		message := "Playlist not found"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Code:    fiber.StatusNotFound,
			Message: message,
		})
	}

	// Update playlist
	podcasts := make([]models.Podcast, 0)
	for _, p := range playlist.Podcasts {
		for _, id := range request.Podcasts {
			if p.ID == id {
				podcasts = append(podcasts, p)
			}
		}
	}

	if len(podcasts) == 0 {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: "No podcasts found for the selected ids",
		})
	}

	err := trx.Association("Podcasts").Replace(podcasts)

	if err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	result = trx.Model(&playlist).Updates(models.Playlist{
		Title:       request.Title,
		Description: request.Description,
	})

	if result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	trx.Commit()
	return c.JSON(definitions.DataResponse[models.PlaylistResource]{
		Code: fiber.StatusOK,
		Data: *playlist.ToResource(),
	})
}

// DeletePlaylist - Delete a playlist
func (plCtrl *PlaylistController) DeletePlaylist(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)

	// Get playlist
	trx := config.DB.Begin()
	playlist := new(models.Playlist)
	result := trx.Where("id = ? AND user_id = ?", c.Params("playlistId"), user.ID).First(&playlist)
	if result.Error != nil {
		message := "Playlist not found"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Code:    fiber.StatusNotFound,
			Message: message,
		})
	}

	// Delete associations
	err := trx.Model(&playlist).Association("Podcasts").Clear()
	if err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Delete playlist
	result = trx.Delete(&playlist)
	if result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	trx.Commit()
	return c.Status(fiber.StatusOK).JSON(definitions.MessageResponse{
		Code:    fiber.StatusOK,
		Message: fmt.Sprintf("Playlist %s deleted successfully", playlist.Title),
	})
}
