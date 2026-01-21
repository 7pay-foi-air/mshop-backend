package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
	"github.com/mshop/account-service/validation"
	"golang.org/x/crypto/bcrypt"
)

// ResetPasswordHandler godoc
// @Summary Reset password using recovery token
// @Description Allows users to reset their password without authentication using a recovery token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/password/reset [post]
func ResetPasswordHandler(c *gin.Context) {
	var req models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	loginRepo := repositories.NewLoginRepository(db.DB)

	user, err := loginRepo.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if user.RecoveryTokenHash == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Recovery token not set"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.RecoveryTokenHash),
		[]byte(req.RecoveryToken),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid recovery token"})
		return
	}

	if !validation.ValidatePassword(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 10 characters long and contain letters and numbers",
		})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	regRepo := repositories.NewRegistrationRepository(db.DB)

	if err := regRepo.ChangePasswordWithRecovery(
		user.UUID,
		string(newHash),
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successful",
	})
}
