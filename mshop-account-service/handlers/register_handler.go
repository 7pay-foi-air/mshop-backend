package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/models"
)

// RegisterHandler godoc
// @Description Receives user and organization registration data together and returns confirmation message
// @Tags Registration
// @Accept json
// @Produce json
// @Param request body models.RegistrationRequest true "Combined registration payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/register [post]
func RegisterHandler(c *gin.Context) {
	var req models.RegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration payload"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration request received",
	})
}
