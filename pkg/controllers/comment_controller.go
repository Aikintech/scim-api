package controllers

import (
	"errors"

	"github.com/aikintech/scim-api/pkg/database"

	"github.com/aikintech/scim-api/pkg/constants"

	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CommentController struct{}

func NewCommentController() *CommentController {
	return &CommentController{}
}

// GetPodcastComments
func (cmtCtrl *CommentController) GetPodcastComments(c *fiber.Ctx) error {
	podcastId := c.Params("podcastId", "")

	// Find podcast
	podcast := new(models.Podcast)
	if result := database.DB.Preload("Comments.User").Where("id = ?", podcastId).Find(&podcast); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
				Message: "No record found",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: result.Error.Error(),
			})
		}
	}

	// Convert to resource
	comments := models.CommentsToResourceCollection(podcast.Comments)

	return c.JSON(comments)
}

// StorePodcastComment - Comment on a podcast
func (cmtCtrl *CommentController) StorePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId", "")

	// Parse request
	request := validation.StorePodcastCommentSchema{}
	if err := c.BodyParser(&request); err != nil {
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

	// Find podcast
	trx := database.DB.Begin()
	podcast := models.Podcast{}
	if result := trx.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Insert comment
	comment := models.Comment{Body: request.Comment, UserID: user.ID, CommentableID: podcastId, CommentableType: podcast.GetPolymorphicType()}
	if result := trx.Model(&models.Comment{}).Create(&comment); result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	trx.Commit()

	comment.User = user

	return c.Status(fiber.StatusCreated).JSON(models.CommentToResource(&comment))
}

// UpdatePodcastComment - Update a podcast comment
func (cmtCtrl *CommentController) UpdatePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId", "")
	commentId := c.Params("commentId", "")

	// Parse request
	request := validation.StorePodcastCommentSchema{}
	if err := c.BodyParser(&request); err != nil {
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

	// Find podcast
	trx := database.DB.Begin()
	podcast := models.Podcast{}
	if result := trx.Model(&models.Podcast{}).Preload("Comments").Where("id = ?", podcastId).First(&podcast); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	// Find comment
	comment, ok := lo.Find(podcast.Comments, func(item *models.Comment) bool {
		return item.ID == commentId
	})
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Message: "Comment not found",
		})
	}

	// Update comment
	if result := trx.Model(&comment).Update("body", request.Comment); result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: result.Error.Error(),
		})
	}

	trx.Commit()

	comment.User = user

	return c.Status(fiber.StatusCreated).JSON(models.CommentToResource(comment))
}

// DeletePodcastComment - Delete a podcast comment
func (cmtCtrl *CommentController) DeletePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(constants.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId", "")
	commentId := c.Params("commentId", "")

	// Find podcast and comment
	trx := database.DB.Begin()
	comment := new(models.Comment)

	if err := trx.Where(&models.Comment{ID: commentId, CommentableID: podcastId, CommentableType: "podcasts", UserID: user.ID}).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
				Message: "No record found",
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
				Message: err.Error(),
			})
		}
	}

	// Delete comment
	if err := trx.Debug().Where("id = ?", commentId).Delete(&comment).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: "Comment deleted successfully",
	})
}
