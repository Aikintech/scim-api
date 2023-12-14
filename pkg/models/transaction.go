package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Transaction struct {
	ID             string    `gorm:"primaryKey;size:40"`
	UserID         string    `gorm:"size:40;index;not null"`
	ReferenceID    string    `gorm:"size:40"`
	Provider       string    `gorm:"size:20"`
	IdempotencyKey string    `gorm:"size:40;index;not null"`
	Currency       string    `gorm:"size:3;not null"`
	Amount         int64     `gorm:"not null"`
	Type           string    `gorm:"size:40;not null"`
	Method         string    `gorm:"size:40;not null"`
	Status         string    `gorm:"size:40;not null;default:'pending'"`
	Description    string    `gorm:"size:255"`
	CreatedAt      time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;not null"`

	User *User `gorm:"foreignKey:UserID"`
}

func (t *Transaction) BeforeCreate(db *gorm.DB) (err error) {
	t.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return
}
