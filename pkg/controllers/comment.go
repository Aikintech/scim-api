package controllers

import (
	"errors"
	"fmt"

	"github.com/aikintech/scim-api/pkg/config"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

func (cmtCtrl *CommentController) GetPodcastComments(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")
	podcast := models.Podcast{}
	result := config.DB.Debug().Preload("Comments.User").Where("id = ?", podcastId).Find(&podcast)

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

	// Convert to resource
	comments := make([]*models.CommentResource, 0)
	for _, c := range podcast.Comments {
		comments = append(comments, c.ToResource())
	}

	return c.JSON(definitions.DataResponse[[]*models.CommentResource]{
		Code: fiber.StatusOK,
		Data: comments,
	})
}

// StorePodcastComment - Comment on a podcast
func (cmtCtrl *CommentController) StorePodcastComment(c *fiber.Ctx) error {
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

	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.CommentResource]{
		Code: fiber.StatusCreated,
		Data: *comment.ToResource(),
	})
}

// UpdatePodcastComment - Update a podcast comment
func (cmtCtrl *CommentController) UpdatePodcastComment(c *fiber.Ctx) error {

	return c.SendString("Like podcast")
}

// DeletePodcastComment - Delete a podcast comment
func (cmtCtrl *CommentController) DeletePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)

	// Find podcast and comment
	trx := config.DB.Begin()
	comment := models.Comment{}
	result := trx.Debug().Where(&models.Comment{
		CommentableID:   c.Params("podcastId"),
		CommentableType: "podcasts",
		UserID:          user.ID,
	}).First(&comment, c.Params("commentId"))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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

	fmt.Println(comment)

	return c.SendString("Like podcast")
}
