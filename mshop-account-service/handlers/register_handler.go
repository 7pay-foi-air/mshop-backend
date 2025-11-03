package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/models"
)

func RegisterHandler(c *gin.Context) {
	var req models.UserRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registration request received",
		"data": gin.H{
			"first_name":   req.FirstName,
			"last_name":    req.LastName,
			"email":        req.Email,
			"phone_number": req.PhoneNumber,
			"username":     req.Username,
			"is_admin":     req.IsAdmin,
		},
	})
}
