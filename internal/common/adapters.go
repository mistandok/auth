package common

import (
	"github.com/mistandok/auth/pkg/user_v1"
)

// RoleNameFromRole get role name for user_v1.Role.
func RoleNameFromRole(role user_v1.Role) string {
	roleName := user_v1.Role_name[int32(role)]

	return roleName
}

// PointerRoleNameFromRole get pointer to role name for user_v1.Role.
func PointerRoleNameFromRole(role *user_v1.Role) *string {
	if role == nil {
		return nil
	}

	roleName := user_v1.Role_name[int32(*role)]

	return &roleName
}

// RoleFromRoleName get user_v1.Role for role name.
func RoleFromRoleName(roleName string) user_v1.Role {
	role := user_v1.Role_value[roleName]

	return user_v1.Role(role)
}
