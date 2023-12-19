package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/utils"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

func (pc *PostController) GetPosts(c *fiber.Ctx) error {
	var total int64
	postType := c.Query("postType", "all")
	search := c.Query("search", "")

	// Query posts
	posts := make([]*models.Post, 0)
	query := database.DB.Model(&models.Post{}).
		Select("posts.*, COUNT(DISTINCT likes.id) AS likesCount, COUNT(DISTINCT comments.id) AS commentsCount").
		Joins("LEFT JOIN likes ON likes.likeable_id = posts.id AND likes.likeable_type = 'posts'").
		Joins("LEFT JOIN comments ON comments.commentable_id = posts.id AND comments.commentable_type = 'posts'").
		Where("posts.title LIKE ?", "%"+search+"%").
		Group("posts.id")

	if postType == "announcement" {
		query = query.Where("posts.is_announcement = ?", true)
	}
	if postType == "blog post" {
		query = query.Where("posts.is_announcement = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := query.Scopes(models.PaginationScope(c)).Preload("User").Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": c.QueryInt("limit", 10),
		"page":  c.QueryInt("page", 1),
		"total": total,
		"items": models.PostsToResourceCollection(posts),
	})
}

func (pc *PostController) GetPost(c *fiber.Ctx) error {
	postId := c.Params("postId")

	// Query post
	post := new(models.Post)
	if err := database.DB.Model(&models.Post{}).Preload("User").First(&post, "id = ?", postId).Error; err != nil {
		status := fiber.StatusBadRequest

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(models.PostToResource(post))
}

// Client handlers
func (pc *PostController) GetPostComments(c *fiber.Ctx) error {
	postId := c.Params("postId")
	isAnnouncement := c.Query("isAnnouncement") == "true"

	// Query post
	post := new(models.Post)
	if err := database.DB.Model(&models.Post{}).
		Preload("Comments").Where("is_announcement = ?", isAnnouncement).
		First(&post, "id = ?", postId).Error; err != nil {
		status := fiber.StatusBadRequest

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(models.CommentsToResourceCollection(post.Comments))
}

func (pc *PostController) CreatePostComment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postId := c.Params("postId")
	request := new(definitions.StoreCommentRequest)

	// Parse request
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

	// Query post
	post := new(models.Post)
	if err := database.DB.Model(&models.Post{}).First(&post, "id = ?", postId).Error; err != nil {
		status := fiber.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}
		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Create comment
	comment := &models.Comment{
		UserID:          user.ID,
		Body:            request.Comment,
		CommentableID:   post.ID,
		CommentableType: post.GetPolymorphicType(),
	}
	if err := database.DB.Create(&comment).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	comment.User = user

	return c.Status(fiber.StatusCreated).JSON(models.CommentToResource(comment))
}

func (pc *PostController) UpdatePostComment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postId := c.Params("postId")
	commentId := c.Params("commentId")
	request := new(definitions.StoreCommentRequest)

	// Parse request
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

	// Query comment
	trx := database.DB.Begin()
	comment := new(models.Comment)
	if err := trx.Model(&models.Comment{}).
		First(&comment, "id = ? AND user_id = ? AND commentable_type = 'posts' AND commentable_id = ?", commentId, user.ID, postId).
		Error; err != nil {
		status := fiber.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}

		trx.Rollback()

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Update comment
	comment.Body = request.Comment
	if err := trx.Save(&comment).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	trx.Commit()

	comment.User = user

	return c.JSON(models.CommentToResource(comment))
}

func (pc *PostController) DeletePostComment(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	postId := c.Params("postId")
	commentId := c.Params("commentId")

	// Delete comment
	trx := database.DB.Begin()
	if err := trx.Where("id = ? AND user_id = ? AND commentable_type = 'posts' AND commentable_id = ?", commentId, user.ID, postId).
		Delete(&models.Comment{}).Error; err != nil {
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

// Backoffice handlers
func (pc *PostController) BackofficeCreatePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	// Parse request
	request := new(definitions.StorePostRequest)
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

	// Create post
	trx := database.DB.Begin()
	post := &models.Post{
		UserID:          user.ID,
		Title:           request.Title,
		Body:            request.Body,
		IsAnnouncement:  request.IsAnnouncement,
		Published:       request.Published,
		Slug:            slug.Make(request.Title) + "-" + time.Now().Format("20060102150405"),
		ExcerptImageURL: request.ExcerptImage,
		MinutesToRead:   request.MinutesToRead,
	}

	if err := trx.Create(&post).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	trx.Commit()

	post.User = user

	return c.JSON(models.PostToResource(post))
}

func (pc *PostController) BackofficeUpdatePost(c *fiber.Ctx) error {
	postId := c.Params("postId")

	// Parse request
	request := new(definitions.StorePostRequest)
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

	// Query post
	trx := database.DB.Begin()
	post := new(models.Post)
	if err := trx.Model(&models.Post{}).First(&post, "id = ?", postId).Error; err != nil {
		status := fiber.StatusBadRequest

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Update post
	post.Title = request.Title
	post.Body = request.Body
	post.Published = request.Published
	post.ExcerptImageURL = request.ExcerptImage
	post.MinutesToRead = request.MinutesToRead
	post.IsAnnouncement = request.IsAnnouncement
	post.Slug = slug.Make(request.Title) + "-" + time.Now().Format("20060102150405")

	if err := trx.Save(&post).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Fetch post with user
	if err := trx.Model(&models.Post{}).Preload("User").First(&post, "id = ?", postId).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	trx.Commit()

	return c.JSON(models.PostToResource(post))
}

func (pc *PostController) BackofficeDeletePost(c *fiber.Ctx) error {
	postId := c.Params("postId")

	// Fetch post
	trx := database.DB.Begin()
	post := new(models.Post)

	if err := trx.Model(&models.Post{}).Where("id = ?", postId).Find(&post).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Delete post
	if err := trx.Delete(&models.Post{}, "id = ?", postId).Error; err != nil {
		trx.Rollback()

		return c.Status(fiber.StatusBadGateway).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	// Delete file from storage
	if len(post.ExcerptImageURL) > 0 {
		go func() {
			if err := utils.DeleteS3File(post.ExcerptImageURL); err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: "Post deleted successfully",
	})
}
