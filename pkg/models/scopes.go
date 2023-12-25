package models

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	Sort       string      `json:"sort"`
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Rows       interface{} `json:"rows"`
}

func PaginationScope(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		offset := (page - 1) * limit

		return db.Offset(offset).Limit(limit)
	}
}
