package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/validation"
)

// LoginHandler godoc
// @Description Receives username and password, returns a success message if payload is valid
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "User login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/login [post]
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login payload"})
		return
	}

	errors := validation.ValidateLogin(req)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, validation.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "error",
			Errors: errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful.",
	})
}
