package models

import (
	"time"

	"github.com/aikintech/scim-api/pkg/constants"
	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Role struct {
	ID          string `gorm:"primaryKey;size:40"`
	Name        string `gorm:"size:40;not null"`
	DisplayName string `gorm:"size:40"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null"`

	// Relations
	Permissions []*Permission `gorm:"many2many:permission_role"`
	Users       []*User       `gorm:"many2many:role_user"`
}

type RoleResource struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	DisplayName string                `json:"displayName"`
	Description string                `json:"description"`
	CreatedAt   time.Time             `json:"createdAt"`
	Permissions []*PermissionResource `json:"permissions"`
}

func (r *Role) BeforeCreate(trx *gorm.DB) error {
	r.ID = nanoid.MustGenerate(constants.ALPHABETS, 30)

	return nil
}

func RoleToResource(r *Role) *RoleResource {
	permissions := PermissionsToResourceCollection(r.Permissions)

	return &RoleResource{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		Permissions: permissions,
	}
}

func RolesToResourceCollection(roles []*Role) []*RoleResource {
	collection := make([]*RoleResource, 0)

	for _, r := range roles {
		collection = append(collection, RoleToResource(r))
	}

	return collection
}
