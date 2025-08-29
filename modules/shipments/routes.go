package shipments

import (
	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/middlewares"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"
)

func Routes(rg *gin.RouterGroup) {
	rg.POST("/", middlewares.JwtAuthMiddleware(), CreateNewShipment)
	rg.GET("/", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin), GetShipmentsList)
	rg.GET("/:id", middlewares.JwtAuthMiddleware(), GetShipmentByID)
	rg.POST("/:id/cancel", middlewares.JwtAuthMiddleware(), CancelShipmentByID)
	rg.POST("/:id/pick-up", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin, roles.RoleCourier), PickupPackageByShipmentID)
	rg.POST("/:id/transit", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin, roles.RoleCourier), TransitPackageByShipmentID)
	rg.POST("/:id/deliver", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin, roles.RoleCourier), DeliverPackageByShipmentID)

	rg.GET("/track/:tracking_number", TrackShipmentHistoriesByTrackingNumber)
}
