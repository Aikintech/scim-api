package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Permission struct {
	ID          string `gorm:"primaryKey;size:40"`
	Name        string `gorm:"size:40;not null"`
	DisplayName string `gorm:"size:40"`
	Module      string `gorm:"size:40"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null"`

	// Relations
	Roles []*Role `gorm:"many2many:permission_role"`
	Users []*User `gorm:"many2many:permission_user"`
}

type PermissionResource struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Module      string `json:"module"`
}

func (p *Permission) BeforeCreate(trx *gorm.DB) error {
	p.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}

func PermissionToResource(p Permission) PermissionResource {
	return PermissionResource{
		ID:          p.ID,
		DisplayName: p.DisplayName,
		Module:      p.Module,
	}
}
