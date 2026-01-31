package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/7pay-foi-air/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/account-service/models"
	"github.com/mshop/account-service/repositories"
)

type UserQuery struct {
	UUIDs []string `form:"uuid"`
}

type UserHandler struct {
	userRepo repositories.UserRepository
}

func NewUserHandler(userRepo repositories.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUsers godoc
// @Summary Get users
// @Description Returns all users or filters by UUIDs if provided
// @Tags User
// @Accept json
// @Produce json
// @Param uuid query []string false "User UUIDs"
// @Success 200 {array} models.UserDB
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	raw := c.Query("uuid")

	parts := []string{}
	if raw != "" {
		parts = strings.Split(raw, ",")
	}

	parsedUUIDs := make([]uuid.UUID, 0, len(parts))
	for _, s := range parts {
		id, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan UUID format."})
			return
		}
		parsedUUIDs = append(parsedUUIDs, id)
	}

	users, err := h.userRepo.GetUsers(parsedUUIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Neuspješno dohvaćanje korisnika",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateProfile godoc
// @Summary Update own profile
// @Description Allows users to update their basic profile information
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.UpdateUserRequest true "User update data"
// @Success 200 {object} models.UserDB
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/profile [patch]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Korisnik nije autentificiran."})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.DateOfBirth != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format datuma. Koristite YYYY-MM-DD."})
			return
		}
		updates["date_of_birth"] = parsedDate
	}
	if req.PhoneNumber != nil {
		if strings.TrimSpace(*req.PhoneNumber) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Broj telefona ne može biti prazan."})
			return
		}
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nema polja za ažuriranje."})
		return
	}

	updatedUser, err := h.userRepo.UpdateUser(uuid.MustParse(claims.UserID), updates)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Broj telefona već postoji."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno ažuriranje profila."})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// UpdateUserByAdmin godoc
// @Summary Update user by admin
// @Description Allows admins to update any user's information including role and status
// @Tags User
// @Accept json
// @Produce json
// @Param userId path string true "User UUID"
// @Param request body models.AdminUpdateUserRequest true "Admin user update data"
// @Success 200 {object} models.UserDB
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/users/{userId} [patch]
func (h *UserHandler) UpdateUserByAdmin(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan ID korisnika."})
		return
	}

	var req models.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	updates := make(map[string]interface{})

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.DateOfBirth != nil {
		parsedDate, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format datuma. Koristite YYYY-MM-DD."})
			return
		}
		updates["date_of_birth"] = parsedDate
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Email != nil {
		if strings.TrimSpace(*req.Email) != "" && !strings.Contains(*req.Email, "@") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format emaila."})
			return
		}
		updates["email"] = *req.Email
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Role != nil {
		validRoles := map[string]bool{"admin": true, "cashier": true}
		if !validRoles[*req.Role] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravna uloga. Mora biti: owner, admin ili cashier."})
			return
		}
		updates["role"] = *req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nema polja za ažuriranje."})
		return
	}

	updatedUser, err := h.userRepo.UpdateUserByAdmin(userID, updates)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Korisnik nije pronađen."})
			return
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email, korisničko ime ili broj telefona već postoji."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno ažuriranje korisnika."})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Soft-delete a user by UUID (sets deleted_at timestamp and is_active to false)
// @Tags User
// @Param uuid path string true "User UUID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/users/{uuid} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("userId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan ID korisnika."})
		return
	}

	rowsAffected, err := h.userRepo.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno brisanje korisnika."})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Korisnik nije pronađen."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Korisnik je uspješno obrisan."})
}
