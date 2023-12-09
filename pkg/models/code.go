package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type VerificationCode struct {
	ID        string    `gorm:"primaryKey;siz:40"`
	UserID    string    `gorm:"not null;index"`
	Code      string    `gorm:"not null;index"`
	Expired   bool      `gorm:"not null;default:false"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	User      *User
}

func (v *VerificationCode) BeforeCreate(tx *gorm.DB) error {
	v.ID = ulid.Make().String()

	return nil
}
