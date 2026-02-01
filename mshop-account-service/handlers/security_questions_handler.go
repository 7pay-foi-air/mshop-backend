package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
	"golang.org/x/crypto/bcrypt"
)

// SetSecurityQuestionsHandler godoc
// @Summary Set security questions and recovery code location
// @Description Saves hashed security questions (answer1+answer2+answer3) and recovery code location
// @Tags Security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SetSecurityQuestionsRequest true "Security questions and recovery location"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/security/questions [post]
func SetSecurityQuestionsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neautorizirani pristup."})
		return
	}

	var req models.SetSecurityQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	combinedAnswers := req.Answer1 + req.Answer2 + req.Answer3
	questionsHash, err := bcrypt.GenerateFromPassword(
		[]byte(combinedAnswers),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška pri kreiranju hash-a."})
		return
	}

	securityRepo := repositories.NewSecurityRepository(db.DB)
	err = securityRepo.SetSecurityQuestions(
		userID.(string),
		string(questionsHash),
		req.RecoveryCodeLocation,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška pri spremanju podataka."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sigurnosna pitanja uspješno spremljena.",
	})
}

// VerifySecurityQuestionsHandler godoc
// @Summary Verify security questions
// @Description Verifies if provided answers match stored hash
// @Tags Security
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SetSecurityQuestionsRequest true "Security answers to verify"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/security/questions/verify [post]
func VerifySecurityQuestionsHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neautorizirani pristup."})
		return
	}

	var req models.SetSecurityQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	securityRepo := repositories.NewSecurityRepository(db.DB)
	storedHash, err := securityRepo.GetSecurityQuestions(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška pri dohvaćanju podataka."})
		return
	}

	combinedAnswers := req.Answer1 + req.Answer2 + req.Answer3
	err = bcrypt.CompareHashAndPassword(
		[]byte(storedHash),
		[]byte(combinedAnswers),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Netočni odgovori.",
			"valid":   false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Odgovori su točni.",
		"valid":   true,
	})
}
