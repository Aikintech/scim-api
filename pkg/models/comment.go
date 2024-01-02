package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Comment struct {
	ID              string    `gorm:"primaryKey;size:40"`
	ParentID        string    `gorm:"index"`
	UserID          string    `gorm:"not null;index"`
	Body            string    `gorm:"not null"`
	CommentableID   string    `gorm:"not null;index"`
	CommentableType string    `gorm:"not null;index"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	User            *User
}

type CommentResource struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	User      *UserRel  `json:"user"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	c.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}

func CommentToResource(c *Comment) CommentResource {
	return CommentResource{
		ID:        c.ID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
		User:      ToUserRelResource(c.User),
	}
}

func CommentsToResourceCollection(comments []*Comment) []CommentResource {
	resources := make([]CommentResource, len(comments))

	for i, comment := range comments {
		resources[i] = CommentToResource(comment)
	}

	return resources
}
