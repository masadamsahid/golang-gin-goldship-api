package shipments

import (
	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/auth"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"
)

func Routes(rg *gin.RouterGroup) {
	rg.POST("/", auth.JwtAuthMiddleware(), CreateNewShipment)
	rg.GET("/", auth.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin), GetShipmentsList)
	rg.GET("/:id", auth.JwtAuthMiddleware(), GetShipmentByID)
	rg.POST("/:id/cancel", auth.JwtAuthMiddleware(), CancelShipmentByID)
	rg.POST("/:id/pick-up", auth.JwtAuthMiddleware(), PickupPackageByShipmentID)
	rg.POST("/:id/transit", auth.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin, roles.RoleCourier), TransitPackageByShipmentID)
	rg.POST("/:id/delivered", auth.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin, roles.RoleCourier), DeliverPackageByShipmentID)
}
