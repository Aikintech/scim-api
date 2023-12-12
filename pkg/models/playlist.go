package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Playlist struct {
	ID          string     `gorm:"primaryKey;size:40"`
	UserID      string     `gorm:"not null;index"`
	Title       string     `gorm:"size:191;not null;index"`
	Description string     `gorm:"size:191"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	Podcasts    []*Podcast `gorm:"many2many:podcast_playlist"`
	User        *User
}

type PlaylistResource struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	CreatedAt   time.Time         `json:"createdAt"`
	Podcasts    []PodcastResource `json:"podcasts"`
}

func (p *Playlist) BeforeCreate(tx *gorm.DB) error {
	p.ID = nanoid.MustGenerate(constants.ALPHABETS, 26)

	return nil
}

func PlaylistToResource(p *Playlist) PlaylistResource {
	return PlaylistResource{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		Podcasts:    PodcastsToResourceCollection(p.Podcasts),
	}
}

func PlaylistsToResourceCollection(playlists []*Playlist) []PlaylistResource {
	resources := make([]PlaylistResource, len(playlists))

	for i, playlist := range playlists {
		resources[i] = PlaylistToResource(playlist)
	}

	return resources
}
