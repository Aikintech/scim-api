package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"

	// "github.com/oklog/ulid/v2"
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
	c.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}
