package shipments

import (
	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/auth"
)

func Routes(rg *gin.RouterGroup) {
	rg.POST("/", auth.JwtAuthMiddleware(), CreateNewShipment)
	rg.GET("/", auth.JwtAuthMiddleware(), CreateNewShipment)
	rg.GET("/:id", auth.JwtAuthMiddleware(), CreateNewShipment)
	rg.POST("/:id/cancel", auth.JwtAuthMiddleware(), CreateNewShipment)
}
