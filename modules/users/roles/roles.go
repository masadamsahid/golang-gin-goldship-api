package roles

type UserRoles = string

const (
	RoleSuperAdmin UserRoles = "SUPERADMIN"
	RoleAdmin      UserRoles = "ADMIN"
	RoleCourier    UserRoles = "COURIER"
	RoleCustomer   UserRoles = "CUSTOMER"
)

var ROLE_LIST []UserRoles = []UserRoles{
	RoleSuperAdmin,
	RoleAdmin,
	RoleCourier,
	RoleCustomer,
}
