package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserToken struct {
	ID          string    `gorm:"primaryKey;size:40"`
	UserID      string    `gorm:"size:40;not null;index"`
	Reference   string    `gorm:"size:40;not null;index"`
	Token       string    `gorm:"not null"`
	Whitelisted bool      `gorm:"not null;index"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	User        *User
}

func (t *UserToken) BeforeCreate(tx *gorm.DB) error {
	t.ID = ulid.Make().String()

	return nil
}
