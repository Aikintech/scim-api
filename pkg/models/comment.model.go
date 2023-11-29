package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Comment struct {
	ID              string `gorm:"primaryKey"`
	Body            string `gorm:"not null"`
	ParentID        string
	UserID          string `gorm:"not null"`
	CommentableID   string
	CommentableType string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	c.ID = ulid.Make().String()

	return nil
}
