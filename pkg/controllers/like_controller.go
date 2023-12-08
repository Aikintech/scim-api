package controllers

import (
	"errors"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LikeController struct{}

func NewLikeController() *LikeController {
	return &LikeController{}
}

// LikePodcast - Like a podcast
func (likeCtrl *LikeController) LikePodcast(c *fiber.Ctx) error {
	// TODO: Optimize this function
	// Fetch podcast
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId")
	podcast := models.Podcast{}
	result := database.DB.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast)
	if result.Error != nil {
		message := "Record not found"
		code := 404

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
			code = fiber.StatusBadRequest
		}

		return c.Status(code).JSON(definitions.MessageResponse{
			Message: message,
		})
	}

	// Fetch like
	like := models.Like{}
	result = database.DB.Model(&models.Like{}).Where(map[string]interface{}{
		"user_id":       user.ID,
		"likeable_type": "podcasts",
		"likeable_id":   podcast.ID,
	}).First(&like)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	}

	// Like or unlike podcast
	message := "Podcast liked successfully"
	if len(like.ID) == 0 {
		result = database.DB.Model(&models.Like{}).Create(&models.Like{
			UserID:       user.ID,
			LikeableID:   podcast.ID,
			LikeableType: "podcasts",
		})

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	} else {
		result = database.DB.Delete(&models.Like{}, "id = ?", like.ID)
		message = "Podcast unliked successfully"

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(definitions.MessageResponse{
		Message: message,
	})
}
