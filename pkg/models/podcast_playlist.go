package models

import (
	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type PodcastPlaylist struct {
	ID         string `gorm:"primaryKey;size:40"`
	PlaylistID string `gorm:"column:playlist_id;index"`
	PodcastID  string `gorm:"column:podcast_id;index"`
	Podcast    Podcast
	Playlist   Playlist
}

func (pp *PodcastPlaylist) BeforeCreate(tx *gorm.DB) error {
	pp.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}

func (PodcastPlaylist) TableName() string {
	return "podcast_playlist"
}
