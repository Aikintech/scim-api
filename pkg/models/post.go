package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Post struct {
	ID              string    `gorm:"primaryKey;size:40"`
	UserID          string    `gorm:"not null"`
	Title           string    `gorm:"not null"`
	Slug            string    `gorm:"not null"`
	Body            string    `gorm:"not null"`
	Published       bool      `gorm:"not null;default:false"`
	ExcerptImageURL string    `gorm:"text"`
	IsAnnouncement  bool      `gorm:"not null;default:false"`
	MinutesToRead   int       `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
	User            User      `gorm:"foreignKey:UserID"`
}

type PostResource struct {
	ID              string           `json:"id"`
	Title           string           `json:"title"`
	Slug            string           `json:"slug"`
	Body            string           `json:"body"`
	Published       bool             `json:"published"`
	ExcerptImageURL string           `json:"excerptImage"`
	IsAnnouncement  bool             `json:"isAnnouncement"`
	MinutesToRead   int              `json:"minutesToRead"`
	CreatedAt       time.Time        `json:"createdAt"`
	User            AuthUserResource `json:"user"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}

func (p *Post) ToResource() PostResource {
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
		User: AuthUserResource{
			ID:            p.User.ID,
			FirstName:     p.User.FirstName,
			LastName:      p.User.LastName,
			Email:         p.User.Email,
			EmailVerified: p.User.EmailVerifiedAt != nil,
			Avatar:        nil,
			Channels:      p.User.Channels,
		},
	}
}
