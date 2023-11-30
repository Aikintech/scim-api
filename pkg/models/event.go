package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Event struct {
	ID              string `gorm:"primaryKey;size:40"`
	Title           string `gorm:"not null"`
	Description     string `gorm:"not null"`
	ExcerptImageURL string
	Location        string    `gorm:"size:255;not null"`
	StartDateTime   time.Time `gorm:"not null"`
	EndDateTime     time.Time
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = ulid.Make().String()

	return nil
}
