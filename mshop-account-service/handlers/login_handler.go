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

	repo := repositories.NewLoginRepository(db.DB)
	user, err := repo.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	orgID := ""
	if user.OrganisationUUID != nil {
		orgID = user.OrganisationUUID.String()
	}
	token, err := auth.GenerateJWT(user.UUID.String(), user.Role, orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user_id": user.UUID,
		"role":    user.Role,
	})
}
