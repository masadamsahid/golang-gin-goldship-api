package users

import (
	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/middlewares"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"
)

func Routes(rg *gin.RouterGroup) {
	rg.GET("/my-shipments", middlewares.JwtAuthMiddleware(), GetMyShipments)
	rg.POST("/:username/change-role", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin), ChangeUserRole)
}
