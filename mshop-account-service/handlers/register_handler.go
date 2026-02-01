package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
	"github.com/mshop/account-service/validation"
	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler godoc
// @Description Receives user and organization registration data together and returns confirmation message
// @Tags Registration
// @Accept json
// @Produce json
// @Param request body models.RegistrationRequest true "Combined registration payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/register [post]
func RegisterHandler(c *gin.Context) {
	var req models.RegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	errors := validation.ValidateRegistration(req)
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, validation.ErrorResponse{
			Code:   http.StatusBadRequest,
			Status: "error",
			Errors: errors,
		})
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno pokretanje transakcije"})
		return
	}
	defer tx.Rollback()

	repo := repositories.NewRegistrationRepository(db.DB)

	userUUID := uuid.New()
	dateOfBirth, _ := time.Parse("2006-01-02", req.DateOfBirth)
	plainPassword := GenerateRandomPassword(10)

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(plainPassword),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno hashiranje lozinke."})
		return
	}

	if err := repo.CreateUser(tx, userUUID, req, string(hashedPassword), dateOfBirth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno kreiranje korisnika."})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno potvrđivanje transakcije."})
		return
	}

	_ = SendRegistrationEmail(
		req.Email,
		req.Username,
		plainPassword,
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registracija uspješna.",
	})
}
