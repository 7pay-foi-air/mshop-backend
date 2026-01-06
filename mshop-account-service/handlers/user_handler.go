package handlers

import (
	"net/http"
	"strings"
	"time"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
		parsedUUIDs = append(parsedUUIDs, id)
	}

	users, err := h.userRepo.GetUsers(parsedUUIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
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
// @Param request body UpdateUserRequest true "User update data"
// @Success 200 {object} models.UserDB
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/profile [patch]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		updates["date_of_birth"] = parsedDate
	}
	if req.PhoneNumber != nil {
		if strings.TrimSpace(*req.PhoneNumber) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number cannot be empty"})
			return
		}
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updatedUser, err := h.userRepo.UpdateUser(userUUID.(uuid.UUID), updates)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
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
// @Param request body AdminUpdateUserRequest true "Admin user update data"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}
		updates["email"] = *req.Email
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Role != nil {
		validRoles := map[string]bool{"owner": true, "admin": true, "cashier": true}
		if !validRoles[*req.Role] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be: owner, admin, or cashier"})
			return
		}
		updates["role"] = *req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	updatedUser, err := h.userRepo.UpdateUserByAdmin(userID, updates)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email, username or phone number already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
