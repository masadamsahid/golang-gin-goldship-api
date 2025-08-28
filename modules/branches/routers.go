package branches

import (
	"github.com/gin-gonic/gin"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers/middlewares"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users/roles"
)

func Routes(rg *gin.RouterGroup) {
	rg.POST("/", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin), HandleCreateBranch)
	rg.GET("/", HandleGetBranchesList)
	rg.GET("/:id", HandleGetBranchByID)
	rg.PUT("/:id", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin), HandleUpdateBranch)
	rg.DELETE("/:id", middlewares.JwtAuthMiddleware(roles.RoleSuperAdmin, roles.RoleAdmin), HandleDeleteBranch)
}
