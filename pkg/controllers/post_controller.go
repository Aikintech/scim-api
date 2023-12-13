package controllers

import (
	"errors"
	"time"

	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/aikintech/scim-api/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// Backoffice handlers
func (pc *PostController) BackofficeGetPosts(c *fiber.Ctx) error {
	var total int64
	isAnnouncement := c.Query("isAnnouncement") == "true"
	search := c.Query("search", "")

	// Query posts
	posts := make([]*models.Post, 0)
	query := database.DB.Model(&models.Post{}).
		Select("posts.*, COUNT(DISTINCT likes.id) AS likesCount, COUNT(DISTINCT comments.id) AS commentsCount").
		Joins("LEFT JOIN likes ON likes.likeable_id = posts.id AND likes.likeable_type = 'posts'").
		Joins("LEFT JOIN comments ON comments.commentable_id = posts.id AND comments.commentable_type = 'posts'").
		Where("posts.is_announcement = ?", isAnnouncement).
		Where("posts.title LIKE ?", "%"+search+"%").
		Group("posts.id")

	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := query.Preload("User").Find(&posts).Error; err != nil {
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

func (pc *PostController) BackofficeCreatePost(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	isAnnouncement := c.Query("isAnnouncement") == "true"

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
		IsAnnouncement:  isAnnouncement,
		Published:       request.Published,
		Slug:            slug.Make(request.Title) + "-" + time.Now().Format("20060102150405"),
		ExcerptImageURL: request.ExcerptImageURL,
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

func (pc *PostController) BackofficeGetPost(c *fiber.Ctx) error {
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
	post.ExcerptImageURL = request.ExcerptImageURL
	post.MinutesToRead = request.MinutesToRead
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

	// Delete post
	trx := database.DB.Begin()
	if err := trx.Delete(&models.Post{}, "id = ?", postId).Error; err != nil {
		status := fiber.StatusBadRequest

		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = fiber.StatusNotFound
		}

		trx.Rollback()

		return c.Status(status).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	trx.Commit()

	return c.JSON(definitions.MessageResponse{
		Message: "Post deleted successfully",
	})
}
