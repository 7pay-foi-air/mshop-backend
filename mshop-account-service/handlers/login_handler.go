package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/auth"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
	"github.com/mshop/account-service/validation"
	"golang.org/x/crypto/bcrypt"
)

// LoginHandler godoc
// @Summary Login user
// @Description Receives username/password, returns access + refresh token
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

	repo := repositories.NewLoginRepository(db.DB)
	user, err := repo.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account is not active"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	orgID := ""
	if user.OrganisationUUID != nil {
		orgID = user.OrganisationUUID.String()
	}

	accessToken, err := auth.GenerateAccessToken(user.UUID.String(), user.Role, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.UUID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"role":          user.Role,
	})
}
