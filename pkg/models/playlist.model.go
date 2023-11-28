package models

import "time"

type Playlist struct {
	ID          string `gorm:"primaryKey"`
	UserID      string
	Title       string     `gorm:"size:191;not null"`
	ShortURL    string     `gorm:"size:191"`
	Description string     `gorm:"size:191"`
	CreateAt    time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
	Podcasts    []*Podcast `gorm:"many2many:podcast_playlist"`
}
