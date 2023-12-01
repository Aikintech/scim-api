package models

import (
	"time"

	"github.com/oklog/ulid/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID              string `gorm:"primaryKey;size:40"`
	ExternalID      string `gorm:"not null"`
	FirstName       string `gorm:"not null"`
	LastName        string `gorm:"not null"`
	Email           string `gorm:"not null"`
	Password        string
	EmailVerifiedAt *time.Time
	SignUpProvider  string `gorm:"not null"`
	Channels        datatypes.JSON
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Playlists       []Playlist
	PrayerRequests  []PrayerRequest
	UserTokens      []UserToken
}

func (model *User) BeforeCreate(*gorm.DB) error {
	model.ID = ulid.Make().String()

	return nil
}

type AuthUserResource struct {
	ID            string         `json:"id"`
	FirstName     string         `json:"firstName"`
	LastName      string         `json:"lastName"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"emailVerified"`
	Avatar        *string        `json:"avatar"`
	Channels      datatypes.JSON `json:"channels"`
}
