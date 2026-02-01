package handlers

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/7pay-foi-air/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/account-service/db"
	"github.com/mshop/account-service/repositories"
)

func GenerateRecoveryToken() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return fmt.Sprintf(
		"%s-%s-%s-%s",
		b[0:4],
		b[4:8],
		b[8:12],
		b[12:16],
	)
}

func GetRecoveryCodeLocationHandler(c *gin.Context) {
	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autentificiran."})
		return
	}

	userUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Neispravan token."})
		return
	}

	securityRepo := repositories.NewSecurityRepository(db.DB)

	location, isActive, isLocked, err := securityRepo.GetRecoveryLocationByUserUUID(userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Korisnik nije pronađen."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Greška pri dohvaćanju podataka."})
		return
	}

	if !isActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Račun nije aktivan."})
		return
	}

	if isLocked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Račun je zaključan."})
		return
	}

	if location == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recovery lokacija nije postavljena."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recovery_code_location": location,
	})
}
