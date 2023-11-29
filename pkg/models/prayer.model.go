package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type PrayerRequest struct {
	ID          string `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Body        string `gorm:"not null"`
	CompletedAt time.Time
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (p *PrayerRequest) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}
