package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	token "github.com/mshop/articles-service/auth"
	"github.com/mshop/articles-service/db"
	"github.com/mshop/articles-service/models"
	"github.com/mshop/articles-service/repositories"
)

type ItemQuery struct {
	UUIDs []string `form:"uuid"`
}

type ItemHandler struct {
	itemRepo repositories.ItemRepository
}

func NewItemHandler(itemRepo repositories.ItemRepository) *ItemHandler {
	return &ItemHandler{
		itemRepo: itemRepo,
	}
}

// GetItems godoc
// @Summary Get items
// @Description Returns all items or filters by UUIDs if provided
// @Tags Items
// @Accept json
// @Produce json
// @Param uuid query []string false "Item UUIDs"
// @Success 200 {array} models.ItemResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/items [get]
func (h *ItemHandler) GetItems(c *gin.Context) {
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

	items, err := h.itemRepo.GetItems(parsedUUIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch items",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

// CreateItem godoc
// @Summary Create item
// @Description Creates a new item with optional image
// @Tags Items
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Item name"
// @Param description formData string true "Item description"
// @Param price formData number true "Item price"
// @Param currency formData string true "Currency code"
// @Param sku formData string false "SKU"
// @Param stock_quantity formData integer true "Stock quantity"
// @Param image formData file false "Item image"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/items [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	currency := c.PostForm("currency")
	sku := c.PostForm("sku")
	stockQuantityStr := c.PostForm("stock_quantity")

	if name == "" || description == "" || priceStr == "" || currency == "" || stockQuantityStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	stockQuantity, err := strconv.ParseInt(stockQuantityStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock quantity format"})
		return
	}

	var imageURL *string
	file, err := c.FormFile("image")
	if err == nil {
		if !isValidImageType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type. Only JPEG, PNG, and GIF are allowed"})
			return
		}

		const maxFileSize = 5 * 1024 * 1024
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image size exceeds 5MB limit"})
			return
		}

		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		uploadPath := filepath.Join("uploads", "items", filename)
		if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		url := fmt.Sprintf("/uploads/items/%s", filename)
		imageURL = &url
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	claims, err := token.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orgID, err := uuid.Parse(claims.OrgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organisation UUID in token"})
		return
	}

	var skuPtr *string
	if sku != "" {
		skuPtr = &sku
	}

	item := models.Item{
		UUIDOrganisation: orgID,
		Name:             name,
		Description:      description,
		Price:            price,
		Currency:         currency,
		SKU:              skuPtr,
		StockQuantity:    int(stockQuantity),
		ImageURL:         imageURL,
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	newID, err := h.itemRepo.CreateItem(item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Item created",
		"uuid_item": newID,
		"image_url": imageURL,
	})
}

// UpdateItem godoc
// @Summary Update item
// @Description Update an item by UUID (accepts JSON data and optional image file, automatically updates the updated_at timestamp)
// @Tags Items
// @Param uuid path string true "Item UUID"
// @Param data formData string true "Item data as JSON string"
// @Param image formData file false "Optional item image"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/items/{uuid} [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	idStr := c.Param("uuid")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	itemData := c.PostForm("data")
	if itemData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing item data"})
		return
	}

	var req models.ItemUpdateRequest
	if err := json.Unmarshal([]byte(itemData), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	var newImageURL *string
	file, err := c.FormFile("image")
	if err == nil {
		if !isValidImageType(file.Header.Get("Content-Type")) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image type. Only JPEG, PNG, and GIF are allowed"})
			return
		}

		const maxFileSize = 5 * 1024 * 1024
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image size exceeds 5MB limit"})
			return
		}

		ext := filepath.Ext(file.Filename)
		filename := uuid.New().String() + ext
		uploadPath := filepath.Join("uploads", "items", filename)

		if err := os.MkdirAll(filepath.Dir(uploadPath), 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}

		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		url := "/uploads/items/" + filename
		newImageURL = &url
	}

	rowsAffected, err := h.itemRepo.UpdateItem(id, req, newImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// DeleteItem godoc
// @Summary Delete item
// @Description Soft-delete an item by UUID (sets deleted_at timestamp)
// @Tags Items
// @Param uuid path string true "Item UUID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/items/{uuid} [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("uuid")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	rowsAffected, err := h.itemRepo.DeleteItem(id)
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

func isValidImageType(contentType string) bool {
	validTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	for _, t := range validTypes {
		if t == contentType {
			return true
		}
	}
	return false
}
