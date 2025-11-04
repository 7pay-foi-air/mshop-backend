package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/docs"
	"github.com/mshop/account-service/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Account Service API
// @version 1.0
// @description API documentation for the account service
// @host localhost:8081
// @BasePath /

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("api/v1/register", handlers.RegisterHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	docs.LaunchSwagger(port)

	r.Run("0.0.0.0:" + port)
}
