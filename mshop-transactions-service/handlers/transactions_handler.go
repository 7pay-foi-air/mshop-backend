package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mshop/auth"
	"github.com/mshop/transactions-service/db"
	"github.com/mshop/transactions-service/models"
	"github.com/mshop/transactions-service/repositories"
)

type TransactionHandler struct {
	repo repositories.TransactionsRepository
}

func NewTransactionHandler(repo repositories.TransactionsRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Creates a purchase transaction with items
// @Tags Transactions
// @Accept json
// @Produce json
//
//	@Param transaction body models.CreateTransactionRequest true "Transaction data"
//
// @Success 201 {object} models.TransactionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req models.CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction must contain items"})
		return
	}

	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(claims.UserID)
	orgID, _ := uuid.Parse(claims.OrgID)

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start DB transaction"})
		return
	}

	txID, err := h.repo.CreateTransaction(tx, req, userID, orgID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	total := 0.0

	for _, item := range req.Items {
		subtotal, err := h.repo.InsertTransactionItem(tx, item, txID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert item"})
			return
		}
		total += subtotal
	}

	if err := h.repo.FinalizeTransaction(tx, txID, total); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize transaction"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB commit failed"})
		return
	}

	c.JSON(http.StatusCreated, models.TransactionResponse{
		UUIDTransaction: txID,
		TotalAmount:     total,
		Currency:        req.Currency,
		IsSuccessful:    true,
	})
}
