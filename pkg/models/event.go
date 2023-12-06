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
	EndDateTime     *time.Time
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

type EventResource struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	ExcerptImageURL string     `json:"excerptImage"`
	Location        string     `json:"location"`
	StartDateTime   time.Time  `json:"startDateTime"`
	EndDateTime     *time.Time `json:"endDateTime"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = ulid.Make().String()

	return nil
}

func (e *Event) ToResource() EventResource {
	return EventResource{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		ExcerptImageURL: e.ExcerptImageURL,
		Location:        e.Location,
		StartDateTime:   e.StartDateTime,
		EndDateTime:     e.EndDateTime,
	}
}
