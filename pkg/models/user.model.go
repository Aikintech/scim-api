package models

import (
	"time"

	"github.com/oklog/ulid/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID              string `gorm:"primaryKey;size:191"`
	ExternalId      string `gorm:"size:191;not null"`
	FirstName       string `gorm:"size:191;not null"`
	LastName        string `gorm:"size:191;not null"`
	Email           string `gorm:"size:191;not null"`
	Password        string `gorm:"size:191"`
	EmailVerifiedAt *time.Time
	SignUpProvider  string `gorm:"size:191;not null"`
	Channels        datatypes.JSON
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Playlists       []Playlist
}

func (model *User) BeforeCreate(*gorm.DB) error {
	model.ID = ulid.Make().String()

	return nil
}
