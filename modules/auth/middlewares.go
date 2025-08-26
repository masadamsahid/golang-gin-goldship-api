package auth

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
)

func JwtAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fmt.Println("Form JWT middleware")
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Auth header not found",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid auth header format",
			})
			return
		}

		strAuthToken := parts[1]
		authToken, err := helpers.VerifyAuthToken(strAuthToken)
		if err != nil || !authToken.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Auth token invalid",
			})
			return
		}

		claims, ok := authToken.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Auth payload invalid",
			})
			return
		}

		// Parse Jwt Payload
		authPayload := helpers.AuthPayload{
			ID:       uint(claims["id"].(float64)),
			Username: claims["username"].(string),
			Email:    claims["email"].(string),
			Role:     claims["role"].(string),
			Exp:      uint(claims["exp"].(float64)),
		}

		// Check Authorization by Role
		isAllowed := slices.Contains(allowedRoles, authPayload.Role)
		// log.Println(authPayload.Role)
		// log.Println(allowedRoles)
		// log.Println(isAllowed)
		if len(allowedRoles) > 0 && !isAllowed {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized to access this resource",
			})
			return
		}

		ctx.Set("user", authPayload)

		ctx.Next()
	}
}
