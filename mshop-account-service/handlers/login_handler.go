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

// LoginHandler godoc
// @Summary Login user
// @Description Receives username/password, returns access + refresh token.
//
//	On first login, generates and returns a recovery token (shown only once).
//
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/login [post]
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
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

	loginRepo := repositories.NewLoginRepository(db.DB)
	user, err := loginRepo.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neispravni podaci za prijavu."})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Račun nije aktivan."})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neispravni podaci za prijavu."})
		return
	}

	orgID := ""
	if user.OrganisationUUID != nil {
		orgID = user.OrganisationUUID.String()
	}

	accessToken, err := auth.GenerateAccessToken(
		user.UUID.String(),
		user.Role,
		orgID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno generiranje pristupnog tokena."})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.UUID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno generiranje osvježavajućeg tokena."})
		return
	}

	response := gin.H{
		"message":       "Prijava uspješna",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"role":          user.Role,
	}

	if user.RecoveryTokenHash == nil {
		recoveryToken := GenerateRecoveryToken()

		recoveryHash, err := bcrypt.GenerateFromPassword(
			[]byte(recoveryToken),
			bcrypt.DefaultCost,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno generiranje recovery tokena."})
			return
		}

		regRepo := repositories.NewRegistrationRepository(db.DB)
		if err := regRepo.SetRecoveryToken(user.UUID, string(recoveryHash)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno pohranjivanje recovery tokena."})
			return
		}

		response["recovery_token"] = recoveryToken
	}

	c.JSON(http.StatusOK, response)
}
