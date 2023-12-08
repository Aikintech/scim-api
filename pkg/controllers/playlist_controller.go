package controllers

import (
	"errors"
	"fmt"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type PlaylistController struct{}

func NewPlaylistController() *PlaylistController {
	return &PlaylistController{}
}

// GetPlaylists - Get user playlists
func (plCtrl *PlaylistController) GetPlaylists(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

	// Get playlists
	playlists := make([]*models.Playlist, 0)
	if result := database.DB.Preload("Podcasts").Where("user_id = ?", user.ID).Find(&playlists); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Convert to resource
	resources := models.PlaylistsToResourceCollection(playlists)

	return c.Status(fiber.StatusOK).JSON(resources)
}

// CreatePlaylist - Create a playlist
func (plCtrl *PlaylistController) CreatePlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

	// Parse body
	request := new(definitions.StorePlaylistRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: "Invalid request body",
		})
	}

	// Validate body
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Message: "Invalid request body",
			Errors:  errs,
		})
	}

	trx := database.DB.Begin()
	podcasts := make([]*models.Podcast, 0)

	// Create playlist
	playlist := models.Playlist{Title: request.Title, Description: request.Description, UserID: user.ID}
	if result := trx.Model(&models.Playlist{}).Create(&playlist); result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Append podcasts to playlist
	if len(request.Podcasts) > 0 {
		if result := trx.Model(&models.Podcast{}).Where(request.Podcasts).Find(&podcasts); result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				trx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Message: result.Error.Error(),
				})
			}
		}

		// Attach podcasts to playlist
		if len(podcasts) > 0 {
			playlistPodcasts := lo.Map(podcasts, func(item *models.Podcast, _ int) models.PodcastPlaylist {
				return models.PodcastPlaylist{PlaylistID: playlist.ID, PodcastID: item.ID}
			})
			if err := trx.Model(&models.PodcastPlaylist{}).Create(&playlistPodcasts).Error; err != nil {
				trx.Rollback()

				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Message: err.Error(),
				})
			}

			playlist.Podcasts = podcasts
		}
	}

	trx.Commit()

	return c.Status(fiber.StatusCreated).JSON(models.PlaylistToResource(&playlist))
}

// GetPlaylist - Get a playlist
func (plCtrl *PlaylistController) GetPlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Get playlist
	playlist := new(models.Playlist)
	if result := database.DB.Preload("Podcasts").Where(models.Playlist{ID: playlistId, UserID: user.ID}).First(&playlist); result.Error != nil {
		message := "Playlist not found"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Message: message,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.PlaylistToResource(playlist))
}

// UpdatePlaylist - updates a user's playlist
func (plCtrl *PlaylistController) UpdatePlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Parse body
	request := new(definitions.StorePlaylistRequest)
	if err := c.BodyParser(request); err != nil {
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

	trx := database.DB.Begin()

	// Get playlist
	playlist := models.Playlist{}
	if result := trx.Preload("Podcasts").Where(&models.Playlist{ID: playlistId, UserID: user.ID}).First(&playlist); result.Error != nil {
		status := 404

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			status = 400
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Update playlist
	playlist.Title = request.Title
	playlist.Description = request.Description
	if err := trx.Save(&playlist).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// // Update playlist podcasts
	// if len(request.Podcasts) > 0 {
	// 	if result := trx.Find(&podcasts, request.Podcasts); result.Error != nil {
	// 		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 			trx.Rollback()

	// 			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
	// 				Message: result.Error.Error(),
	// 			})
	// 		}
	// 	}

	// 	// Sync podcasts to playlist
	// 	if len(podcasts) > 0 {
	// 		// Delete associations
	// 		if err := trx.Model(&playlist).Association("Podcasts").Clear(); err != nil {
	// 			trx.Rollback()

	// 			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
	// 				Message: err.Error(),
	// 			})
	// 		}

	// 		// Update playlist podcasts
	// 		playlistPodcasts := lo.Map(podcasts, func(item *models.Podcast, _ int) models.PodcastPlaylist {
	// 			return models.PodcastPlaylist{
	// 				PlaylistID: playlistId,
	// 				PodcastID:  item.ID,
	// 			}
	// 		})
	// 		if err := trx.Model(&models.PodcastPlaylist{}).Create(&playlistPodcasts).Error; err != nil {
	// 			trx.Rollback()

	// 			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
	// 				Message: err.Error(),
	// 			})
	// 		}
	// 	}
	// }

	trx.Commit()

	return c.JSON(models.PlaylistToResource(&playlist))
}

// DeletePlaylist - Delete a playlist
func (plCtrl *PlaylistController) DeletePlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

	// Get playlist
	trx := database.DB.Begin()
	playlist := new(models.Playlist)
	result := trx.Where("id = ? AND user_id = ?", c.Params("playlistId"), user.ID).First(&playlist)
	if result.Error != nil {
		message := "Playlist not found"

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
		}

		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Message: message,
		})
	}

	// Delete associations
	err := trx.Model(&playlist).Association("Podcasts").Clear()
	if err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Delete playlist
	result = trx.Delete(&playlist)
	if result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	trx.Commit()
	return c.Status(fiber.StatusOK).JSON(definitions.MessageResponse{
		Message: fmt.Sprintf("Playlist %s deleted successfully", playlist.Title),
	})
}

// DeletePlaylistPodcasts - Removes podcasts from a playlist
func (plCtrl *PlaylistController) DeletePlaylistPodcasts(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Parse request
	request := definitions.PlaylistPodcastsRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(&request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Get playlist with its podcasts
	trx := database.DB.Begin()
	playlist := new(models.Playlist)
	if err := trx.Preload("Podcasts").Where(&models.Playlist{ID: playlistId, UserID: user.ID}).First(&playlist).Error; err != nil {
		status := fiber.StatusNotFound

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusBadRequest
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Delete playlist podcasts
	intersection := lo.Intersect(request.Podcasts, lo.Map(playlist.Podcasts, func(item *models.Podcast, _ int) string {
		return item.ID
	}))
	toBeDeleted := lo.Map(intersection, func(item string, _ int) models.PodcastPlaylist {
		return models.PodcastPlaylist{
			PlaylistID: playlistId,
			PodcastID:  item,
		}
	})

	if len(toBeDeleted) > 0 {
		if err := trx.Delete(&models.PodcastPlaylist{}, toBeDeleted).Error; err != nil {
			trx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: "OK",
	})
}

// AddPlaylistPodcasts
func (plCtrl *PlaylistController) AddPlaylistPodcasts(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Parse request
	request := new(definitions.PlaylistPodcastsRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := validation.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&definitions.ValidationErrsResponse{
			Errors: errs,
		})
	}

	// Find playlist
	trx := database.DB.Begin()
	playlist := new(models.Playlist)
	if err := trx.Preload("Podcasts").Where(&models.Playlist{ID: playlistId, UserID: user.ID}).First(&playlist).Error; err != nil {
		status := fiber.StatusNotFound

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusBadRequest
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Find podcasts
	mappedPodcasts := lo.Map(playlist.Podcasts, func(item *models.Podcast, _ int) string {
		return item.ID
	})
	diff := lo.Filter(request.Podcasts, func(item string, index int) bool {
		return !lo.Contains(mappedPodcasts, item)
	})

	if len(diff) > 0 {
		podcasts := make([]models.Podcast, len(diff))
		if err := trx.Debug().Find(&podcasts, diff).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}

		// Attach playlist podcasts
		playlistPodcasts := lo.Map(podcasts, func(item models.Podcast, _ int) models.PodcastPlaylist {
			return models.PodcastPlaylist{PlaylistID: playlistId, PodcastID: item.ID}
		})
		if err := trx.Model(&models.PodcastPlaylist{}).Create(&playlistPodcasts).Error; err != nil {
			trx.Rollback()

			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	trx.Commit()

	return c.Status(fiber.StatusCreated).JSON(definitions.MessageResponse{
		Message: "Podcasts added to playlist successfully",
	})
}
