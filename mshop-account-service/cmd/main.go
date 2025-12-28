package main

import (
	"log"
	"net/http"
	"os"

	"github.com/7pay-foi-air/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/docs"
	"github.com/mshop/account-service/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Account Service API
// @version 1.0
// @description API documentation for the account service
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	secret := os.Getenv("JWT_SECRET")
	auth.SetAccesSecretKey(secret)
	auth.SetRefreshSecretKey(secret)

	db.Init()

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("api/v1/login", handlers.LoginHandler)

	r.POST("/api/v1/refresh", handlers.RefreshTokenHandler)

	protected := r.Group("/api/v1")
	protected.Use(auth.RequiredAuth())

	protected.POST("/password/change", handlers.ChangePasswordHandler)

	admin := protected.Group("/")
	admin.Use(auth.RequiredAdmin())

	admin.POST("/register", handlers.RegisterHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	docs.LaunchSwagger(port)

	r.Run("0.0.0.0:" + port)
}
