package main

import (
	"log"
	"net/http"
	"os"

	"github.com/7pay-foi-air/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mshop/articles-service/db"
	"github.com/mshop/articles-service/docs"
	"github.com/mshop/articles-service/handlers"
	"github.com/mshop/articles-service/repositories"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Articles Service API
// @version 1.0
// @description API documentation for the articles service
// @host localhost:8082
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

	itemRepo := repositories.NewItemRepository(db.DB)
	itemHandler := handlers.NewItemHandler(itemRepo)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.Static("/uploads", "./uploads")
	protected := r.Group("/api/v1")
	protected.Use(auth.RequiredAuth())

	admin := protected.Group("/")
	admin.Use(auth.RequiredAdmin())

	protected.GET("/items", itemHandler.GetItems)
	admin.POST("/items", itemHandler.CreateItem)
	admin.DELETE("/items/:uuid", itemHandler.DeleteItem)
	admin.PUT("/items/:uuid", itemHandler.UpdateItem)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	docs.LaunchSwagger(port)

	r.Run("0.0.0.0:" + port)
}
