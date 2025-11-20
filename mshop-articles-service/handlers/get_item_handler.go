package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/articles-service/db"
	"github.com/mshop/articles-service/repositories"
)

type ItemQuery struct {
	UUIDs []string `form:"uuid"`
}

// GetItemsHandler godoc
// @Summary Get items
// @Description Returns all items or filters by UUIDs if provided
// @Tags Items
// @Accept json
// @Produce json
// @Param uuid query []string false "Item UUIDs"
// @Success 200 {array} models.ItemResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/items [get]
func GetItemsHandler(c *gin.Context) {
	var query ItemQuery

	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	parsedUUIDs := make([]uuid.UUID, 0, len(query.UUIDs))
	for _, s := range query.UUIDs {
		id, err := uuid.Parse(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
			return
		}
		parsedUUIDs = append(parsedUUIDs, id)
	}

	repo := repositories.GetItemRepository(db.DB)
	items, err := repo.GetItems(parsedUUIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch items",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}
