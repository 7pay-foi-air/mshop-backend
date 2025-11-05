package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
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

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	orgUUID := uuid.New()
	if err := repositories.CreateOrganisation(tx, orgUUID, req.Organization); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organisation"})
		return
	}

	userUUID := uuid.New()
	dateOfBirth, _ := time.Parse("2006-01-02", req.User.DateOfBirth)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)

	if err := repositories.CreateUser(tx, userUUID, orgUUID, req.User, string(hashedPassword), dateOfBirth); err != nil {
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
