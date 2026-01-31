package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/7pay-foi-air/auth"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan zahtjev."})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno pokretanje DB transakcije."})
		return
	}

	txID, err := h.repo.CreateTransaction(tx, req, userID, orgID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno kreiranje transakcije."})
		return
	}

	total := 0.0

	for _, item := range req.Items {
		subtotal, err := h.repo.InsertTransactionItem(tx, item, txID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno umetanje stavke."})
			return
		}
		total += subtotal
	}

	if len(req.Items) == 0 {
		if req.TotalAmount != nil {
			total = *req.TotalAmount
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Prazna transakcija mora uključivati totalAmount."})
			return
		}
	}

	if err := h.repo.FinalizeTransaction(tx, txID, total); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno finaliziranje transakcije."})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješan DB commit."})
		return
	}

	c.JSON(http.StatusCreated, models.TransactionResponse{
		UUIDTransaction: txID,
		TotalAmount:     total,
		Currency:        req.Currency,
		IsSuccessful:    true,
	})
}

// GetUserTransactions godoc
// @Summary Get user transactions
// @Description Retrieves successful transactions for the authenticated user. Admins and owners can see all organization transactions. Supports optional date range filtering.
// @Tags Transactions
// @Produce json
// @Param start_date query string false "Start date (YYYY-MM-DD format)" example(2024-01-01)
// @Param end_date query string false "End date (YYYY-MM-DD format, inclusive)" example(2024-12-31)
// @Success 200 {object} models.TransactionHistoryResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID := uuid.MustParse(claims.UserID)
	orgID := uuid.MustParse(claims.OrgID)
	role := claims.Role

	isAdmin := role == "admin" || role == "owner"

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	transactions, err := h.repo.GetUserTransactions(userID, orgID, isAdmin, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransactionDetails godoc
// @Summary Get transaction details
// @Description Returns transaction basic info + items (items can be empty)
// @Tags Transactions
// @Produce json
// @Param id path string true "Transaction UUID"
// @Success 200 {object} models.TransactionDetailsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetTransactionDetails(c *gin.Context) {
	idStr := c.Param("id")
	txID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Neispravan format UUID-a."})
		return
	}

	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	orgID := uuid.MustParse(claims.OrgID)
	userID := uuid.MustParse(claims.UserID)
	role := claims.Role
	isAdmin := role == "admin" || role == "owner"

	header, err := h.repo.GetTransactionByID(txID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transakcija nije pronađena."})
		return
	}

	if header.UUIDOrganisation != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Zabranjeno"})
		return
	}

	if !isAdmin && header.UUIDUser != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Zabranjeno"})
		return
	}

	items, err := h.repo.GetTransactionItems(txID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Neuspješno dohvaćanje stavki transakcije."})
		return
	}

	c.JSON(http.StatusOK, models.TransactionDetailsResponse{
		UUIDTransaction:     header.UUIDTransaction,
		TransactionType:     header.TransactionType,
		TotalAmount:         header.TotalAmount,
		Currency:            header.Currency,
		TransactionDate:     header.TransactionDate,
		Items:               items,
		TransactionRefundID: header.TransactionRefundID,
		PaymentMethod:       header.PaymentMethod,
	})
}
