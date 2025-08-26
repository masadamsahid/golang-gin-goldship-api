package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/helpers"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users"
	"golang.org/x/crypto/bcrypt"
)

func HandleRegister(ctx *gin.Context) {
	var body RegisteUserDto
	err := ctx.ShouldBind(&body)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}

		log.Printf("%+v\n", validationErrors)
		errs := helpers.HandleValidationErrors(validationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": errs})
		return
	}

	hashedPwd, err := helpers.HashPassword(body.Password)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	sqlCreateNewUser := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, role, created_at, updated_at`

	var newUser users.User
	err = db.DB.QueryRow(sqlCreateNewUser, body.Username, body.Email, hashedPwd).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.Role,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{
				"message": "'username' or 'email' already taken",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed registering new user",
		})
		return
	}

	authToken, err := helpers.CreateAuthToken(helpers.AuthTokenClaims{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Role:     newUser.Role,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success register",
		"data":    authToken,
	})
}

func HandleLogin(ctx *gin.Context) {
	var body LoginUserDto
	err := ctx.ShouldBind(&body)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}

		log.Printf("%+v\n", validationErrors)
		errs := helpers.HandleValidationErrors(validationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed", "errors": errs})
		return
	}

	sqlGetUserByUsername := `SELECT id, username, email, role, "password" FROM users WHERE username = $1 OR email = $1`

	var user users.User
	err = db.DB.QueryRow(sqlGetUserByUsername, body.UsernameEmail).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Role,
		&user.Password,
	)
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Wrong credentials",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed logging in",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		if err.Error() == bcrypt.ErrMismatchedHashAndPassword.Error() {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Wrong credentials",
			})
			return
		}

		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	authToken, err := helpers.CreateAuthToken(helpers.AuthTokenClaims{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed creating token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success logging in",
		"data":    authToken,
	})
}
