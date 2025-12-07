package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mshop/auth"
	"github.com/mshop/transactions-service/db"
	"github.com/mshop/transactions-service/docs"
	"github.com/mshop/transactions-service/handlers"
	"github.com/mshop/transactions-service/repositories"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Transactions Service API
// @version 1.0
// @description API documentation for the transactions service
// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @name Authorization
// @in header

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db.Init()

	secret := os.Getenv("JWT_SECRET")
	auth.SetAccesSecretKey(secret)
	auth.SetRefreshSecretKey(secret)

	repo := repositories.NewTransactionsRepository(db.DB)
	handler := handlers.NewTransactionHandler(repo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	protected := r.Group("/api/v1")
	protected.Use(auth.RequiredAuth())

	admin := protected.Group("/")
	admin.Use(auth.RequiredAdmin())

	protected.POST("/transactions", handler.CreateTransaction)
	protected.GET("/transactions", handler.GetUserTransactions)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	docs.LaunchSwagger(port)

	r.Run("0.0.0.0:" + port)
}
