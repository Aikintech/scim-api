package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Post struct {
	ID             string    `gorm:"primaryKey"`
	Title          string    `gorm:"not null"`
	Slug           string    `gorm:"not null"`
	Body           string    `gorm:"not null"`
	Published      bool      `gorm:"not null;default:false"`
	ExcerptImage   string    `gorm:"text"`
	IsAnnouncement bool      `gorm:"not null;default:false"`
	MinutesToRead  int       `gorm:"not null;default:0"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}
