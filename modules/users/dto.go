package users

import "github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"

type ChangeUserRoleDto struct {
	Role roles.UserRoles `json:"role" binding:"required"`
}
