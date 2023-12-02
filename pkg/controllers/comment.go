package controllers

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/aikintech/scim/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPodcastComments(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")
	podcast := models.Podcast{}
	result := config.DB.Debug().Preload("Comments").Where("id = ?", podcastId).Find(&podcast)

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

	return c.JSON(fiber.Map{
		"data": podcast.Comments,
	})
}

// StorePodcastComment - Comment on a podcast
func StorePodcastComment(c *fiber.Ctx) error {
	// Parse request
	request := new(validation.StorePodcastCommentSchema)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	errs := utils.ValidateStruct(request)
	if len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Find podcast
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId", "")
	podcast := models.Podcast{}
	result := config.DB.Where("id = ?", podcastId).First(&podcast)
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

	comment := models.Comment{
		Body:            request.Comment,
		UserID:          user.ID,
		CommentableID:   podcast.ID,
		CommentableType: "podcasts",
	}
	result = config.DB.Model(&models.Comment{}).Create(&comment)
	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.Comment]{
		Code: fiber.StatusCreated,
		Data: comment,
	})
}

// UpdatePodcastComment - Update a podcast comment
func UpdatePodcastComment(c *fiber.Ctx) error {

	return c.SendString("Like podcast")
}
