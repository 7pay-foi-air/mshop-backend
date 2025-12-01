package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mshop/articles-service/db"
	"github.com/mshop/articles-service/docs"
	"github.com/mshop/articles-service/handlers"
	"github.com/mshop/articles-service/repositories"
	"github.com/mshop/auth"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Articles Service API
// @version 1.0
// @description API documentation for the articles service
// @host localhost:8082
// @BasePath /

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

	r.GET("/api/v1/items", itemHandler.GetItems)
	r.POST("/api/v1/items", itemHandler.CreateItem)
	r.DELETE("/api/v1/items/:uuid", itemHandler.DeleteItem)
	r.PUT("/api/v1/items/:uuid", itemHandler.UpdateItem)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	docs.LaunchSwagger(port)

	r.Run("0.0.0.0:" + port)
}
