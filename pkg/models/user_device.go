package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type UserDevice struct {
	ID         string    `gorm:"primaryKey"`
	UserID     string    `gorm:"not null;index"`
	DeviceOS   string    `gorm:"not null"` // android, ios, web
	DeviceType string    `gorm:"not null"` // phone, tablet, desktop
	FCMToken   string    `gorm:"not null"`
	Active     bool      `gorm:"not null;default:true"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`

	// TODO: Check for user device when logging in
}

func (ud *UserDevice) BeforeCreate(trx *gorm.DB) error {
	ud.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}
