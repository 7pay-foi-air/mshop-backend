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
// @Router /api/v1/register [post]
func RegisterHandler(c *gin.Context) {
	var req models.RegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration payload"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	repo := repositories.NewRegistrationRepository(db.DB)

	orgUUID := uuid.New()

	if err := repo.CreateOrganisation(tx, orgUUID, req.Organization); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organisation"})
		return
	}

	userUUID := uuid.New()
	dateOfBirth, _ := time.Parse("2006-01-02", req.User.DateOfBirth)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)

	if err := repo.CreateUser(tx, userUUID, orgUUID, req.User, string(hashedPassword), dateOfBirth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successfull",
	})
}
