package controllers

import (
	"github.com/aikintech/scim/pkg/config"
	"github.com/aikintech/scim/pkg/definitions"
	"github.com/aikintech/scim/pkg/models"
	"github.com/aikintech/scim/pkg/utils"
	"github.com/aikintech/scim/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

func ClientGetPodcastComments(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")
	podcast := models.Podcast{}
	result := config.DB.Debug().Preload("Comments").Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast)
	println()
	println()
	println()
	println()

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

// ClientStorePodcastComment - Comment on a podcast
func ClientStorePodcastComment(c *fiber.Ctx) error {
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
		UserID:          ulid.Make().String(),
		CommentableID:   podcast.ID,
		CommentableType: "Podcast",
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

// ClientUpdatePodcastComment - Update a podcast comment
func ClientUpdatePodcastComment(c *fiber.Ctx) error {

	return c.SendString("Like podcast")
}
