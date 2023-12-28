package database

import (
	"fmt"

	"github.com/aikintech/scim-api/pkg/models"
	"github.com/samber/lo"
)

func SeedRoles() {
	// Fetch all permissions
	trx := DB.Begin()
	permissions := make([]*models.Permission, 0)
	if err := trx.Model(&models.Permission{}).Find(&permissions).Error; err != nil {
		fmt.Println(err.Error())
	}

	// Upsert role
	role := new(models.Role)
	if err := trx.Model(&models.Role{}).
		Preload("Permissions").
		Assign(models.Role{Name: "super-admin", DisplayName: "Super Administrator", Description: ""}).
		Where("name = ?", "super-admin").
		FirstOrCreate(&role).Error; err != nil {
		trx.Rollback()
		fmt.Println(err.Error())
	}

	mappedPermissions := lo.Map(permissions, func(item *models.Permission, index int) string {
		return item.ID
	})
	mappedRolePermissions := lo.Map(role.Permissions, func(item *models.Permission, index int) string {
		return item.ID
	})

	// Assign permissions
	_, differences := lo.Difference(mappedRolePermissions, mappedPermissions)
	permissionRole := lo.Map(differences, func(item string, index int) models.PermissionRole {
		return models.PermissionRole{RoleID: role.ID, PermissionID: item}
	})

	if err := trx.Model(&models.PermissionRole{}).CreateInBatches(permissionRole, 100).Error; err != nil {
		trx.Rollback()
		fmt.Println(err.Error())
	}

	trx.Commit()
}
