package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xendit/xendit-go/v7/invoice"

	"github.com/joho/godotenv"
	"github.com/masadamsahid/golang-gin-goldship-api/db"
	xenditService "github.com/masadamsahid/golang-gin-goldship-api/helpers/xendit-service"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/auth"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/branches"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/shipments"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/webhooks"
)

func main() {
	xenditService.InitXendit()

	log.Println(
		string(invoice.INVOICESTATUS_PENDING),
		invoice.INVOICESTATUS_PAID,
		invoice.INVOICESTATUS_SETTLED,
		invoice.INVOICESTATUS_EXPIRED,
		invoice.INVOICESTATUS_XENDIT_ENUM_DEFAULT_FALLBACK,
	)

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
	shipments.Routes(api.Group("/shipments"))
	webhooks.Routes(api.Group("/webhooks"))

	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "8080"
	}
	r.Run(":" + PORT)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables from OS")
	} else {
		log.Println("Loading .env success")
	}
}
