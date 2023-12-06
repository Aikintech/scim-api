package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Comment struct {
	ID              string `gorm:"primaryKey;size:40"`
	Body            string `gorm:"not null"`
	ParentID        string
	UserID          string `gorm:"not null"`
	CommentableID   string
	CommentableType string
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	User            *User
}

type CommentResource struct {
	ID        string            `json:"id"`
	Body      string            `json:"body"`
	CreatedAt time.Time         `json:"createdAt"`
	User      *AuthUserResource `json:"user"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	c.ID = ulid.Make().String()

	return nil
}

func (c *Comment) ToResource() *CommentResource {
	return &CommentResource{
		ID:        c.ID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		User:      c.User.ToResource(),
	}
}
