package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Podcast struct {
	ID          string     `gorm:"primaryKey"`
	GUID        string     `gorm:"size:191;not null"`
	Author      string     `gorm:"size:191;not null"`
	Title       string     `gorm:"size:191;not null"`
	SubTitle    string     `gorm:"size:191;not null"`
	Summary     string     `gorm:"size:191;not null"`
	Description string     `gorm:"type:LONGTEXT"`
	Duration    string     `gorm:"size:191;not null"`
	Image       string     `gorm:"size:191;not null"`
	Url         string     `gorm:"type:TEXT;not null"`
	Published   bool       `gorm:"type:TINYINT;not null;default:true"`
	PublishedAt *time.Time `gorm:"not null"`
	Meta        datatypes.JSON
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (p *Podcast) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}

type PodcastResource struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Duration    string    `json:"duration"`
	Image       string    `json:"image"`
	Url         string    `json:"url"`
	Published   bool      `json:"published"`
	PublishedAt time.Time `json:"publishedAt"`
}
