package handlers

import (
	"net/http"

	"github.com/7pay-foi-air/auth"
	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
	"github.com/mshop/account-service/validation"
	"golang.org/x/crypto/bcrypt"
)

// ChangePasswordHandler godoc
// @Summary Change password using recovery token
// @Description Changes user password by validating the recovery token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/password/change [post]
func ChangePasswordHandler(c *gin.Context) {
	value, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neautoriziran pristup"})
		return
	}

	claims := value.(*auth.Claims)
	userID := claims.UserID

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	loginRepo := repositories.NewLoginRepository(db.DB)
	user, err := loginRepo.GetUserByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije pronađen."})
		return
	}

	if user.RecoveryTokenHash == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Recovery token nije postavljen."})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.RecoveryTokenHash),
		[]byte(req.RecoveryToken),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neispravan recovery token."})
		return
	}

	if !validation.ValidatePassword(req.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lozinka mora imati najmanje 10 znakova i sadržavati slova i brojeve.",
		})
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno hashiranje lozinke."})
		return
	}

	regRepo := repositories.NewRegistrationRepository(db.DB)
	if err := regRepo.ChangePasswordWithRecovery(
		user.UUID,
		string(newPasswordHash),
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno ažuriranje lozinke."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lozinka je uspješno promijenjena.",
	})
}
