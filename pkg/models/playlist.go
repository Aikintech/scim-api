package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Playlist struct {
	ID          string `gorm:"primaryKey;size:40"`
	UserID      string
	Title       string    `gorm:"size:191;not null"`
	ShortURL    *string   `gorm:"size:191"`
	Description string    `gorm:"size:191"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	Podcasts    []Podcast `gorm:"many2many:podcast_playlist"`
}

type PlaylistResource struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	ShortURL    *string            `json:"shortUrl"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"createdAt"`
	Podcasts    []*PodcastResource `json:"podcasts"`
}

func (p *Playlist) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}

func (p *Playlist) ToResource() *PlaylistResource {
	podcasts := make([]*PodcastResource, len(p.Podcasts))
	for i, podcast := range p.Podcasts {
		podcasts[i] = podcast.ToResource()
	}

	return &PlaylistResource{
		ID:          p.ID,
		Title:       p.Title,
		ShortURL:    p.ShortURL,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		Podcasts:    podcasts,
	}
}
