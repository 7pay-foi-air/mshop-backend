package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/7pay-foi-air/auth"
	"github.com/mshop/transactions-service/db"
	"github.com/mshop/transactions-service/models"
)

// RefundTransaction godoc
// @Summary Refund transaction
// @Description Creates a refund transaction for an existing purchase
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RefundTransactionRequest true "Refund data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/transactions/refund [post]
func (h *TransactionHandler) RefundTransaction(c *gin.Context) {
	var req models.RefundTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
		return
	}

	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID := uuid.MustParse(claims.UserID)

	original, err := h.repo.GetTransactionByID(req.UUIDTransaction)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transakcija nije pronađena."})
		return
	}

	if original.TransactionRefundID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transakcija je već refundirana."})
		return
	}

	description := req.Description
	if description == "" {
		description = "Refund transaction"
	}

	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno pokretanje transakcije."})
		return
	}
	defer tx.Rollback()

	refundID, err := h.repo.CreateRefundTransaction(tx, original, userID, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Neuspješno kreiranje refund transakcije.",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno potvrđivanje transakcije."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"refund_transaction_id": refundID,
	})
}
