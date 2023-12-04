package controllers

import (
	"errors"

	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
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
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId")
	podcast := models.Podcast{}
	result := config.DB.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast)
	if result.Error != nil {
		message := "Record not found"
		code := 404

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			message = result.Error.Error()
			code = fiber.StatusBadRequest
		}

		return c.Status(code).JSON(definitions.MessageResponse{
			Code:    code,
			Message: message,
		})
	}

	// Fetch like
	like := models.Like{}
	result = config.DB.Model(&models.Like{}).Where(map[string]interface{}{
		"user_id":       user.ID,
		"likeable_type": "podcasts",
		"likeable_id":   podcast.ID,
	}).First(&like)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	// Like or unlike podcast
	message := "Podcast liked successfully"
	if len(like.ID) == 0 {
		result = config.DB.Model(&models.Like{}).Create(&models.Like{
			UserID:       user.ID,
			LikeableID:   podcast.ID,
			LikeableType: "podcasts",
		})

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	} else {
		result = config.DB.Delete(&models.Like{}, "id = ?", like.ID)
		message = "Podcast unliked successfully"

		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Code:    fiber.StatusBadRequest,
				Message: result.Error.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(definitions.MessageResponse{
		Code:    fiber.StatusOK,
		Message: message,
	})
}
