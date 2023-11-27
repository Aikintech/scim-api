package models

import (
	"strconv"

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

func PaginateScope(c *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		offset := (page - 1) * limit

		return db.Offset(offset).Limit(limit)
	}
}
