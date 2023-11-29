package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Like struct {
	ID           string `gorm:"primaryKey"`
	UserID       string `gorm:"not null"`
	LikeableID   string
	LikeableType string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *Like) BeforeCreate(tx *gorm.DB) error {
	c.ID = ulid.Make().String()

	return nil
}
