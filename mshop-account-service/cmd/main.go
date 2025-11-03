package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/handlers"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/register", handlers.RegisterHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run("0.0.0.0:" + port)
}
