package users

import (
	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/middlewares"
)

func Routes(rg *gin.RouterGroup) {
	rg.GET("/my-shipments", middlewares.JwtAuthMiddleware(), GetMyShipments)
}
