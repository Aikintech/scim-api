package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type PrayerRequest struct {
	ID          string `gorm:"primaryKey;size:40"`
	UserID      string `gorm:"not null"`
	Title       string `gorm:"not null"`
	Body        string `gorm:"not null"`
	CompletedAt *time.Time
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	User        *User
}

type PrayerRequestResource struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"description"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
}

func (p *PrayerRequest) BeforeCreate(tx *gorm.DB) error {
	p.ID = ulid.Make().String()

	return nil
}

func (p *PrayerRequest) ToResource() PrayerRequestResource {
	return PrayerRequestResource{
		ID:          p.ID,
		Title:       p.Title,
		Body:        p.Body,
		CompletedAt: p.CompletedAt,
		CreatedAt:   p.CreatedAt,
	}
}
