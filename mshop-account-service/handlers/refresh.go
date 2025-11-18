package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/auth"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/repositories"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenHandler
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body handlers.RefreshRequest true "Refresh token payload"
// @Success 200 {object} map[string]string
// @Router /api/v1/refresh [post]
func RefreshTokenHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	userID, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	repo := repositories.NewLoginRepository(db.DB)
	user, _ := repo.GetUserByID(userID)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	orgID := ""
	if user.OrganisationUUID != nil {
		orgID = user.OrganisationUUID.String()
	}

	newAccessToken, err := auth.GenerateAccessToken(user.UUID.String(), user.Role, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
