package database

import (
	"fmt"

	"github.com/aikintech/scim-api/pkg/constants"
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

	// Upsert roles
	regularUser := new(models.Role)
	if err := trx.Model(&models.Role{}).
		Assign(models.Role{Name: constants.REGULAR_USER_ROLE, DisplayName: "Regular User", Description: ""}).
		Where("name = ?", constants.REGULAR_USER_ROLE).
		FirstOrCreate(&regularUser).Error; err != nil {
		trx.Rollback()
		fmt.Println(err.Error())
	}
	trx.Commit()

	trx = DB.Begin()
	superAdmin := new(models.Role)
	if err := trx.Model(&models.Role{}).
		Preload("Permissions").
		Assign(models.Role{Name: constants.SUPER_ADMIN_USER_ROLE, DisplayName: "Super Administrator", Description: ""}).
		Where("name = ?", constants.SUPER_ADMIN_USER_ROLE).
		FirstOrCreate(&superAdmin).Error; err != nil {
		trx.Rollback()
		fmt.Println(err.Error())
	}

	mappedPermissions := lo.Map(permissions, func(item *models.Permission, index int) string {
		return item.ID
	})
	mappedRolePermissions := lo.Map(superAdmin.Permissions, func(item *models.Permission, index int) string {
		return item.ID
	})

	// Assign super admin permissions
	_, differences := lo.Difference(mappedRolePermissions, mappedPermissions)
	permissionRole := lo.Map(differences, func(item string, index int) models.PermissionRole {
		return models.PermissionRole{RoleID: superAdmin.ID, PermissionID: item}
	})

	if err := trx.Model(&models.PermissionRole{}).CreateInBatches(permissionRole, 100).Error; err != nil {
		trx.Rollback()
		fmt.Println(err.Error())
	}

	trx.Commit()
}
