package models

import (
	"github.com/oklog/ulid/v2"
	"time"

	"gorm.io/datatypes"
)

// User model
type User struct {
	ID              string         `json:"id" gorm:"primaryKey;size:191"`
	FirstName       string         `json:"firstName" gorm:"size:191;not null"`
	LastName        string         `json:"lastName" gorm:"size:191;not null"`
	Email           string         `json:"email" gorm:"size:191;not null"`
	Password        string         `json:"password" gorm:"size:191"`
	EmailVerifiedAt *time.Time     `json:"emailVerifiedAt"`
	SignUpProvider  string         `json:"signUpProvider" gorm:"size:191;not null"`
	Channels        datatypes.JSON `json:"channels"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

func (model *User) BeforeCreate() {
	model.ID = ulid.Make().String()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	return
}

func (model *User) BeforeUpdate() {
	model.UpdatedAt = time.Now()

	return
}
