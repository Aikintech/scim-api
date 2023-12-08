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
	trx := database.DB.Begin()

	// Fetch podcast
	podcast := new(models.Podcast)
	if result := trx.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Fetch like
	like := new(models.Like)
	if result := trx.Model(&models.Like{}).Where(models.Like{UserID: user.ID, LikeableType: podcast.GetPolymorphicType(), LikeableID: podcast.ID}).First(&like); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	}

	// Like or unlike podcast
	message := "Podcast liked successfully"
	if len(like.ID) == 0 {
		if result := trx.Model(&models.Like{}).Create(&models.Like{UserID: user.ID, LikeableID: podcast.ID, LikeableType: podcast.GetPolymorphicType()}); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	} else {
		if result := trx.Delete(&models.Like{}, "id = ?", like.ID); result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}

		message = "Podcast un-liked successfully"
	}

	trx.Commit()

	return c.Status(fiber.StatusOK).JSON(definitions.MessageResponse{
		Message: message,
	})
}
