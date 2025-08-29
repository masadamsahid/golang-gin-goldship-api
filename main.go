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
	"github.com/masadamsahid/golang-gin-goldship-api/modules/users"
	"github.com/masadamsahid/golang-gin-goldship-api/modules/webhooks"

	scalargo "github.com/bdpiprava/scalar-go"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Routes under "/api"
	auth.Routes(api.Group("/auth"))
	branches.Routes(api.Group("/branches"))
	shipments.Routes(api.Group("/shipments"))
	users.Routes(api.Group("/users"))
	webhooks.Routes(api.Group("/webhooks"))

	// OpenAPI
	openAPIYAMLDocs, err := os.ReadFile("./docs/openapi-specs.json")
	if err != nil {
		log.Println("Error loading OpenAPI .yaml file:", err)
	}

	r.GET("/openapi", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/plain", openAPIYAMLDocs)
	})

	// Scalar UI
	html, err := scalargo.NewV2(
		scalargo.WithSpecURL("/openapi"),
		scalargo.WithTheme(scalargo.ThemeBluePlanet),
		scalargo.WithMetaDataOpts(
			scalargo.WithTitle("🚚 Goldship Logistic API (Golang + Gin)"),
		),
	)
	if err != nil {
		log.Println("Error Scalar-Go:", err)
	}

	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", []byte(html))
	})

	// Gin-Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/openapi"),
	))

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
