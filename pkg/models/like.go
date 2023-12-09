package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Like struct {
	ID           string    `gorm:"primaryKey;size:40"`
	UserID       string    `gorm:"not null;index"`
	LikeableID   string    `gorm:"not null;index"`
	LikeableType string    `gorm:"not null;index"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	User         *User
}

type LikeResource struct{}

func (c *Like) BeforeCreate(tx *gorm.DB) error {
	c.ID = ulid.Make().String()

	return nil
}
