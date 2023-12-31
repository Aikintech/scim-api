package controllers

import (
	"strings"

	"github.com/aikintech/scim-api/pkg/database"
	"github.com/aikintech/scim-api/pkg/definitions"
	"github.com/aikintech/scim-api/pkg/models"
	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// Backoffice routes
func (u *UserController) BackofficeGetUsers(c *fiber.Ctx) error {
	var total int64
	search := strings.TrimSpace(c.Query("search", ""))

	// Get users
	users := make([]*models.User, 0)
	query := database.DB.Model(&models.User{}).
		Where("first_name LIKE ?", "%"+search+"%").
		Or("last_name LIKE ?", "%"+search+"%").
		Or("email LIKE ?", "%"+search+"%")

	if err := query.Count(&total).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	if err := query.Scopes(models.PaginationScope(c)).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(definitions.MessageResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(definitions.Map{
		"limit": c.QueryInt("limit", 10),
		"page":  c.QueryInt("page", 1),
		"total": total,
		"items": models.UsersToResourceCollection(users),
	})
}
