package models

import (
	"strings"
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type PrayerRequest struct {
	ID          string `gorm:"primaryKey;size:40"`
	UserID      string `gorm:"not null;index"`
	Title       string `gorm:"not null"`
	Body        string `gorm:"not null"`
	PhoneNumber string `gorm:"not null;index"`
	Status      string `gorm:"not null;index;default:pending"`
	CompletedAt *time.Time
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	User        *User
}

type PrayerRequestResource struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Body        string     `json:"description"`
	PhoneNumber string     `json:"phoneNumber"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	User        UserRel    `json:"user"`
}

func (p *PrayerRequest) BeforeCreate(tx *gorm.DB) error {
	p.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}

func PrayerToResource(p *PrayerRequest) PrayerRequestResource {
	return PrayerRequestResource{
		ID:          p.ID,
		Title:       strings.TrimSpace(p.Title),
		Body:        strings.TrimSpace(p.Body),
		PhoneNumber: p.PhoneNumber,
		Status:      p.Status,
		CompletedAt: p.CompletedAt,
		CreatedAt:   p.CreatedAt,
		User:        ToUserRelResource(p.User),
	}
}

func PrayersToResourceCollection(prayers []*PrayerRequest) []PrayerRequestResource {
	resources := make([]PrayerRequestResource, len(prayers))

	for i, prayer := range prayers {
		resources[i] = PrayerToResource(prayer)
	}

	return resources
}
