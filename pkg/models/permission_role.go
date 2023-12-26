package models

type PermissionRole struct {
	PermissionID string
	RoleID       string
}

func (pr *PermissionRole) TableName() string {
	return "permission_role"
}
