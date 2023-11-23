package models

import (
	"time"

	"gorm.io/datatypes"
)

// Model
type User struct {
	ID              string         `json:"id" gorm:"primaryKey;size:191"`
	FirstName       string         `json:"firstName" gorm:"size:191;not null"`
	LastName        string         `json:"lastName" gorm:"size:191;not null"`
	Email           string         `json:"email" gorm:"size:191;not null"`
	Password        string         `json:"password" gorm:"size:191"`
	EmailVerifiedAt time.Time      `json:"emailVerifiedAt"`
	SignUpProvider  string         `json:"signUpProvider" gorm:"size:191;not null"`
	Channels        datatypes.JSON `json:"channels"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}
