package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
