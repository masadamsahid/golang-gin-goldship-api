package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// var randJWTSecret = rand.Text()
// var randJWTSecretArrOfByte = []byte(randJWTSecret)

var jwtSecret string
var jwtSecretArrOfByte []byte

func InitJWT() {
	jwtSecret = os.Getenv("JWT_SECRET_KEY")
	jwtSecretArrOfByte = []byte(jwtSecret)

	// log.Println("JWT secret:", jwtSecret)
}

type AuthPayload struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Exp      uint   `json:"exp"`
}

type AuthTokenClaims struct {
	ID       uint
	Username string
	Email    string
	Role     string
}

func CreateAuthToken(claims AuthTokenClaims) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       claims.ID,
		"username": claims.Username,
		"email":    claims.Email,
		"role":     claims.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	// log.Println("Token claim", tokenClaims)

	authToken, err := tokenClaims.SignedString(jwtSecretArrOfByte)
	if err != nil {
		log.Println("Error creating auth token:", err)
		return "", err
	}

	return authToken, nil
}

func VerifyAuthToken(strAuthToken string) (*jwt.Token, error) {
	authToken, err := jwt.Parse(strAuthToken, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return jwtSecretArrOfByte, nil
	})

	// log.Println(err)

	return authToken, err
}

func ParseJWTUserFromCtx(ctx *gin.Context, userContainer *AuthPayload) error {
	u, ok := ctx.Get("user")
	if !ok {
		return errors.New("failed to get user from context")
	}

	*userContainer, ok = u.(AuthPayload)
	if !ok {
		return errors.New("failed convert user from context to AuthPayload")
	}

	return nil
}
