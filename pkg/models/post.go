package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Post struct {
	ID              string    `gorm:"primaryKey;size:40"`
	UserID          string    `gorm:"not null;index"`
	Title           string    `gorm:"not null;index"`
	Slug            string    `gorm:"not null;index"`
	Body            string    `gorm:"not null"`
	Published       bool      `gorm:"not null;default:false"`
	ExcerptImageURL string    `gorm:"text"`
	IsAnnouncement  bool      `gorm:"not null;default:false"`
	MinutesToRead   float32   `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	User            *User
	Comments        []*Comment `gorm:"polymorphic:Commentable"`
	Likes           []*Like    `gorm:"polymorphic:Likeable"`
}

type PostResource struct {
	ID              string           `json:"id"`
	Title           string           `json:"title"`
	Slug            string           `json:"slug"`
	Body            string           `json:"body"`
	Published       bool             `json:"published"`
	ExcerptImageURL string           `json:"excerptImage"`
	IsAnnouncement  bool             `json:"isAnnouncement"`
	MinutesToRead   float32          `json:"minutesToRead"`
	CreatedAt       time.Time        `json:"createdAt"`
	User            AuthUserResource `json:"user"`
	LikesCount      int              `json:"likesCount"`
	CommentsCount   int              `json:"commentsCount"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = nanoid.MustGenerate(constants.ALPHABETS, 26)

	return nil
}

func (p *Post) GetPolymorphicType() string {
	return "posts"
}

func PostToResource(p *Post) PostResource {
	return PostResource{
		ID:              p.ID,
		Title:           p.Title,
		Slug:            p.Slug,
		Body:            p.Body,
		Published:       p.Published,
		ExcerptImageURL: p.ExcerptImageURL,
		IsAnnouncement:  p.IsAnnouncement,
		MinutesToRead:   p.MinutesToRead,
		CreatedAt:       p.CreatedAt,
		User:            UserToResource(p.User),
	}
}

func PostsToResourceCollection(posts []*Post) []PostResource {
	postResources := make([]PostResource, len(posts))
	for i, post := range posts {
		postResources[i] = PostToResource(post)
	}

	return postResources
}
