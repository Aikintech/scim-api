package models

import (
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

const ALPHABETS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type UserDevice struct {
	ID         string    `gorm:"primaryKey"`
	UserID     string    `gorm:"not null"`
	DeviceID   string    `gorm:"not null"`
	FCMToken   string    `gorm:"not null"`
	DeviceType string    `gorm:"not null"` // android, ios, web
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

func (ud *UserDevice) BeforeCreate(trx *gorm.DB) error {
	ud.ID = nanoid.MustGenerate(ALPHABETS, 21)

	return nil
}
