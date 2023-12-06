package models

import (
	"time"

	"github.com/oklog/ulid/v2"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User model
type User struct {
	ID                string `gorm:"primaryKey;size:40"`
	ExternalID        string `gorm:"not null"`
	FirstName         string `gorm:"not null"`
	LastName          string `gorm:"not null"`
	Email             string `gorm:"not null"`
	Password          string
	EmailVerifiedAt   *time.Time
	SignUpProvider    string `gorm:"not null"`
	Avatar            *string
	Channels          datatypes.JSON
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Playlists         []*Playlist
	PrayerRequests    []*PrayerRequest
	UserTokens        []*UserToken
	VerificationCodes []*VerificationCode
	Posts             []*Post
	Comments          []*Comment
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

func (u *User) BeforeCreate(*gorm.DB) error {
	u.ID = ulid.Make().String()

	return nil
}

func (u *User) ToResource() *AuthUserResource {
	// Generate avatarURL
	return &AuthUserResource{
		ID:            u.ID,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Email:         u.Email,
		EmailVerified: u.EmailVerifiedAt != nil,
		Avatar:        u.Avatar,
		Channels:      u.Channels,
	}
}
