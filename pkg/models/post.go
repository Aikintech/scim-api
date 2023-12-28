package models

import (
	"fmt"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/utils"
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
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Body            string    `json:"body"`
	Published       bool      `json:"published"`
	ExcerptImage    string    `json:"excerptImage"`
	ExcerptImageKey string    `json:"excerptImageKey"`
	IsAnnouncement  bool      `json:"isAnnouncement"`
	MinutesToRead   float32   `json:"minutesToRead"`
	CreatedAt       time.Time `json:"createdAt"`
	User            UserRel   `json:"user"`
	LikesCount      int       `json:"likesCount"`
	CommentsCount   int       `json:"commentsCount"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = nanoid.MustGenerate(constants.ALPHABETS, 26)

	return nil
}

func (p *Post) GetPolymorphicType() string {
	return "posts"
}

func PostToResource(p *Post) PostResource {
	excerpt, err := utils.GenerateS3FileURL(p.ExcerptImageURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	return PostResource{
		ID:              p.ID,
		Title:           p.Title,
		Slug:            p.Slug,
		Body:            p.Body,
		Published:       p.Published,
		ExcerptImage:    excerpt,
		ExcerptImageKey: p.ExcerptImageURL,
		IsAnnouncement:  p.IsAnnouncement,
		MinutesToRead:   p.MinutesToRead,
		CreatedAt:       p.CreatedAt,
		User:            ToUserRelResource(p.User),
	}
}

func PostsToResourceCollection(posts []*Post) []PostResource {
	postResources := make([]PostResource, len(posts))
	for i, post := range posts {
		postResources[i] = PostToResource(post)
	}

	return postResources
}
