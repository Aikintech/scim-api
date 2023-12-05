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
	podcast := models.Podcast{}
	if result := config.DB.Preload("Comments.User").Where("id = ?", podcastId).Find(&podcast); result.Error != nil {
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
		user := c.User
		avatar, _ := utils.GenerateS3FileURL(*user.Avatar)

		comment := (models.CommentResource{
			ID:        c.ID,
			Body:      c.Body,
			CreatedAt: c.CreatedAt,
			User: &models.AuthUserResource{
				ID:            c.User.ID,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Email:         user.Email,
				EmailVerified: user.EmailVerifiedAt != nil,
				Avatar:        &avatar,
				Channels:      user.Channels,
			},
		})

		comments = append(comments, &comment)
	}

	return c.JSON(definitions.DataResponse[[]*models.CommentResource]{
		Code: fiber.StatusOK,
		Data: comments,
	})
}

// StorePodcastComment - Comment on a podcast
func (cmtCtrl *CommentController) StorePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId", "")

	// Parse request
	request := validation.StorePodcastCommentSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Find podcast
	trx := config.DB.Begin()
	podcast := models.Podcast{}
	if result := trx.Model(&models.Podcast{}).Where("id = ?", podcastId).First(&podcast); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	// Insert comment
	comment := models.Comment{Body: request.Comment, UserID: user.ID, CommentableID: podcastId, CommentableType: podcast.GetPolymorphicType()}
	if result := trx.Debug().Model(&models.Comment{}).Create(&comment); result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	trx.Commit()

	avatar, _ := utils.GenerateS3FileURL(*user.Avatar)

	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.CommentResource]{
		Code: fiber.StatusCreated,
		Data: models.CommentResource{
			ID:        comment.ID,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			User: &models.AuthUserResource{
				ID:            user.ID,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Email:         user.Email,
				EmailVerified: user.EmailVerifiedAt != nil,
				Avatar:        &avatar,
				Channels:      user.Channels,
			},
		},
	})
}

// UpdatePodcastComment - Update a podcast comment
func (cmtCtrl *CommentController) UpdatePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)
	podcastId := c.Params("podcastId", "")
	commentId := c.Params("commentId", "")

	// Parse request
	request := validation.StorePodcastCommentSchema{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Validate request
	if errs := utils.ValidateStruct(request); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(definitions.ValidationErrsResponse{
			Code:   fiber.StatusUnprocessableEntity,
			Errors: errs,
		})
	}

	// Find podcast
	trx := config.DB.Begin()
	podcast := models.Podcast{}
	if result := trx.Model(&models.Podcast{}).Preload("Comments").Where("id = ?", podcastId).First(&podcast); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	// Find comment
	comment, ok := lo.Find(podcast.Comments, func(item *models.Comment) bool {
		return item.ID == commentId
	})
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(definitions.MessageResponse{
			Code:    fiber.StatusNotFound,
			Message: "Comment not found",
		})
	}

	// Update comment
	if result := trx.Model(&comment).Update("body", request.Comment); result.Error != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Code:    fiber.StatusBadRequest,
			Message: result.Error.Error(),
		})
	}

	trx.Commit()

	avatar, _ := utils.GenerateS3FileURL(*user.Avatar)

	return c.Status(fiber.StatusCreated).JSON(definitions.DataResponse[models.CommentResource]{
		Code: fiber.StatusCreated,
		Data: models.CommentResource{
			ID:        comment.ID,
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			User: &models.AuthUserResource{
				ID:            user.ID,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Email:         user.Email,
				EmailVerified: user.EmailVerifiedAt != nil,
				Avatar:        &avatar,
				Channels:      user.Channels,
			},
		},
	})
}

// DeletePodcastComment - Delete a podcast comment
func (cmtCtrl *CommentController) DeletePodcastComment(c *fiber.Ctx) error {
	user := c.Locals(config.USER_CONTEXT_KEY).(*models.User)

	// Find podcast and comment
	trx := config.DB.Begin()
	comment := models.Comment{}
	result := trx.Where(&models.Comment{
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
