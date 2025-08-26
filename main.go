package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/auth"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/branches"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from OS")
	}

	defer db.StopDB()
	db.ConnectDB()

	r := gin.Default()
	r.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, world! Server is ok",
		})
	})

	api := r.Group("/api")
	auth.Routes(api.Group("/auth"))
	branches.Routes(api.Group("/branches"))

	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "8080"
	}
	r.Run(":" + PORT)
}
