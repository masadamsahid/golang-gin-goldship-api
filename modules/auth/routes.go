package auth

import (
	"github.com/gin-gonic/gin"
)

func Routes(rg *gin.RouterGroup) {
	rg.POST("/register", HandleRegister)
	rg.POST("/login", HandleLogin)
}
