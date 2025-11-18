package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mshop/articles-service/db"
	"github.com/mshop/articles-service/models"
	"github.com/mshop/articles-service/repositories"
)

// CreateItem godoc
// @Summary Create new item
// @Description Adds a new item for an organisation
// @Tags Items
// @Accept json
// @Produce json
// @Param request body models.ItemCreateRequest true "Item data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/items [post]
func CreateItem(c *gin.Context) {
	var req models.ItemCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	repo := repositories.NewItemRepository(db.DB)
	id, err := repo.CreateItem(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Item created", "id": id})
}
