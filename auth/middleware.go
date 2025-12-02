package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequiredAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ValidateAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)

		c.Next()
	}
}

func RequiredAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			c.Abort()
			return
		}

		claims := value.(*Claims)

		if claims.Role != "admin" && claims.Role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
