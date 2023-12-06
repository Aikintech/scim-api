package controllers

import (
	"errors"
	"fmt"
	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aws/smithy-go/ptr"
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
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

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

	trx := config.DB.Begin()
	podcasts := make([]*models.Podcast, 0)

	// Create playlist
	playlist := models.Playlist{Title: request.Title, Description: request.Description, UserID: user.ID}
	if result := trx.Model(&models.Playlist{}).Create(&playlist); result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	// Append podcasts to playlist
	if len(request.Podcasts) > 0 {
		if result := trx.Debug().Model(&models.Podcast{}).Where(request.Podcasts).Find(&podcasts); result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				trx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Code:    fiber.StatusBadRequest,
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
					Code:    fiber.StatusBadRequest,
					Message: err.Error(),
				})
			}
		}
	}

	trx.Commit()
	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.PlaylistResource]{
		Code:    fiber.StatusCreated,
		Message: ptr.String("Playlist created successfully"),
		Data: models.PlaylistResource{
			ID:          playlist.ID,
			Title:       playlist.Title,
			ShortURL:    nil,
			Description: playlist.Description,
			CreatedAt:   playlist.CreatedAt,
			Podcasts: lo.Map(podcasts, func(item *models.Podcast, index int) *models.PodcastResource {
				return item.ToResource()
			}),
		},
	})
}

// GetPlaylist - Get a playlist
func (plCtrl *PlaylistController) GetPlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId")

	// Get playlist
	playlist := new(models.Playlist)
	if result := config.DB.Preload("Podcasts").Where(models.Playlist{ID: podcastId, UserID: user.ID}).First(&playlist); result.Error != nil {
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

// UpdatePlaylist - updates a user's playlist
func (plCtrl *PlaylistController) UpdatePlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Parse body
	request := new(definitions.StorePlaylistRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	trx := config.DB.Begin()
	podcasts := make([]models.Podcast, 0)

	// Get playlist
	playlist := models.Playlist{}
	if result := trx.Where(&models.Playlist{ID: playlistId, UserID: user.ID}).First(&playlist); result.Error != nil {
		status := 404

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			status = 400
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Code:    status,
			Message: result.Error.Error(),
		})
	}

	// Update playlist
	playlist.Title = request.Title
	playlist.Description = request.Description
	if err := trx.Save(&playlist).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Update playlist podcasts
	if len(request.Podcasts) > 0 {
		if result := trx.Find(&podcasts, request.Podcasts); result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				trx.Rollback()

				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Code:    fiber.StatusBadRequest,
					Message: result.Error.Error(),
				})
			}
		}

		// Sync podcasts to playlist
		if len(podcasts) > 0 {
			// Delete associations
			if err := trx.Model(&playlist).Association("Podcasts").Clear(); err != nil {
				trx.Rollback()

				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Code:    fiber.StatusBadRequest,
					Message: err.Error(),
				})
			}

			// Update playlist podcasts
			playlistPodcasts := lo.Map(podcasts, func(item models.Podcast, _ int) models.PodcastPlaylist {
				return models.PodcastPlaylist{
					PlaylistID: playlistId,
					PodcastID:  item.ID,
				}
			})
			if err := trx.Model(&models.PodcastPlaylist{}).Create(&playlistPodcasts).Error; err != nil {
				trx.Rollback()

				return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
					Code:    fiber.StatusBadRequest,
					Message: err.Error(),
				})
			}
		}
	}

	trx.Commit()

	return c.JSON(definitions.DataResponse[models.PlaylistResource]{
		Code: fiber.StatusOK,
		Data: models.PlaylistResource{
			ID:          playlist.ID,
			Title:       playlist.Title,
			Description: playlist.Description,
			CreatedAt:   playlist.CreatedAt,
			Podcasts: lo.Map(podcasts, func(item models.Podcast, _ int) *models.PodcastResource {
				return item.ToResource()
			}),
		},
	})
}

// DeletePlaylist - Delete a playlist
func (plCtrl *PlaylistController) DeletePlaylist(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)

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

// DeletePlaylistPodcasts - Removes podcasts from a playlist
func (plCtrl *PlaylistController) DeletePlaylistPodcasts(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Parse request
	request := definitions.PlaylistPodcastsRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Get playlist with its podcasts
	trx := config.DB.Begin()
	playlist := new(models.Playlist)
	if result := trx.Preload("Podcasts").Where(models.Playlist{ID: playlistId, UserID: user.ID}).Find(&playlist); result != nil {
		status := 404

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			status = 400
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Code:    status,
			Message: result.Error.Error(),
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
	if err := trx.Delete(&models.PodcastPlaylist{}, toBeDeleted).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Code:    fiber.StatusOK,
		Message: "OK",
	})
}

// AddPlaylistPodcasts
func (plCtrl *PlaylistController) AddPlaylistPodcasts(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	playlistId := c.Params("playlistId")

	// Parse request
	request := definitions.PlaylistPodcastsRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Find playlist
	trx := config.DB.Begin()
	playlist := new(models.Playlist)
	if err := trx.Where(&models.Playlist{ID: playlistId, UserID: user.ID}).First(&playlist).Error; err != nil {
		status := fiber.StatusNotFound

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusBadRequest
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Code:    status,
			Message: err.Error(),
		})
	}

	// Find podcasts
	podcasts := make([]models.Podcast, 0)
	if err := trx.Find(&podcasts, request.Podcasts).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
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
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	trx.Commit()

	return c.Status(fiber.StatusCreated).JSON(definitions.MessageResponse{
		Code:    fiber.StatusCreated,
		Message: "Podcasts added to playlist successfully",
	})
}
