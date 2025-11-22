package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/articles-service/db"
	"github.com/mshop/articles-service/repositories"
)

// DeleteItemHandler godoc
// @Summary Delete item
// @Description Soft-delete an item by UUID (sets deleted_at timestamp)
// @Tags Items
// @Param uuid path string true "Item UUID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/items/{uuid} [delete]
func DeleteItemHandler(c *gin.Context) {
	idStr := c.Param("uuid")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	rowsAffected, err := repositories.DeleteItemRepository(db.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
