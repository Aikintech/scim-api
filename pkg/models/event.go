package models

import (
	"fmt"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/utils"
	naoid "github.com/matoous/go-nanoid/v2"
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
	Published       bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

type EventResource struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	ExcerptImageURL *string    `json:"excerptImage"`
	Location        string     `json:"location"`
	StartDateTime   time.Time  `json:"startDateTime"`
	EndDateTime     *time.Time `json:"endDateTime"`
	Published       bool       `json:"published"`
	CreatedAt       time.Time  `json:"createdAt"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	e.ID = naoid.MustGenerate(constants.ALPHABETS, 26)

	return nil
}

func (e *Event) ToResource() EventResource {
	// Generate avatarURL
	excerptImage, err := utils.GenerateS3FileURL(e.ExcerptImageURL)
	if err != nil {
		fmt.Println("Error generating avatar url", err.Error())
	}

	return EventResource{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		ExcerptImageURL: &excerptImage,
		Location:        e.Location,
		StartDateTime:   e.StartDateTime,
		EndDateTime:     e.EndDateTime,
		Published:       e.Published,
		CreatedAt:       e.CreatedAt,
	}
}

func EventsToResourceCollection(events []Event) []EventResource {
	resources := make([]EventResource, len(events))

	for i, event := range events {
		resources[i] = event.ToResource()
	}

	return resources
}
