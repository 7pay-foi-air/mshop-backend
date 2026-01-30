package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/7pay-foi-air/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mshop/transactions-service/models"
	"github.com/mshop/transactions-service/repositories"
)

type ReportHandler struct {
	repo repositories.TransactionsRepository
}

func NewReportHandler(repo repositories.TransactionsRepository) *ReportHandler {
	return &ReportHandler{repo: repo}
}

// helper funkcija koja formatira float64 u hrvatski currency
func toHrCurrency(amount float64) string {
	intPart := int64(amount)
	fracPart := int64((amount - float64(intPart)) * 100)

	// tisućni separator
	intStr := fmt.Sprintf("%d", intPart)
	var result string
	for i, c := range reverseString(intStr) {
		if i > 0 && i%3 == 0 {
			result = "." + result
		}
		result = string(c) + result
	}

	return fmt.Sprintf("%s,%02d", result, fracPart)
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CreateReport generira CSV report i šalje ga na email (koristi repo.GetUserTransactions)
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req models.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// validate dates (format YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date (expected YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date (expected YYYY-MM-DD)"})
		return
	}

	// auth
	claims, err := auth.GetTokenClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userID := uuid.MustParse(claims.UserID)
	orgID := uuid.MustParse(claims.OrgID)
	role := claims.Role
	isAdmin := role == "admin" || role == "owner"

	// fetch transactions
	raw, err := h.repo.GetUserTransactions(userID, orgID, isAdmin, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions: " + err.Error()})
		return
	}

	rows := append(raw.SuccessfulTransactions, raw.RefundedTransactions...)

	sort.Slice(rows, func(i, j int) bool {
		ti, err1 := time.Parse(time.RFC3339, rows[i].TransactionDate)
		tj, err2 := time.Parse(time.RFC3339, rows[j].TransactionDate)
		if err1 != nil || err2 != nil {
			return false
		}
		return ti.After(tj)
	})

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{
		"Iznos",
		"Valuta",
		"Datum i vrijeme",
		"Tip",
		"Metoda placanja",
	}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}

	for _, r := range rows {
		formattedDate := r.TransactionDate
		if t, err := time.Parse(time.RFC3339, r.TransactionDate); err == nil {
			formattedDate = t.Format("02.01.2006 15:04")
		}

		method := r.PaymentMethod
		switch method {
		case "card_payment":
			method = "Kartica"
		case "crypto":
			method = "Kriptovalute"
		default:
			method = "Nepoznato"
		}

		txType := r.TransactionType
		switch txType {
		case "Purchase":
			txType = "Kupnja"
		case "Refund":
			txType = "Povrat"
		default:
			txType = ""
		}

		record := []string{
			toHrCurrency(r.TotalAmount),
			r.Currency,
			formattedDate,
			txType,
			method,
		}
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV record"})
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize CSV: " + err.Error()})
		return
	}

	csvBytes := buf.Bytes()
	filename := fmt.Sprintf("transactions_%s_%s.csv", req.StartDate, req.EndDate)
	attachment := &EmailAttachment{
		Filename: filename,
		Content:  csvBytes,
		MimeType: "text/csv",
	}

	start, _ := time.Parse("2006-01-02", req.StartDate)
	end, _ := time.Parse("2006-01-02", req.EndDate)

	subject := fmt.Sprintf(
		"[mShop] Transakcijski izvještaj %s - %s",
		start.Format("02.01.2006"),
		end.Format("02.01.2006"),
	)
	body := "Poštovani,\n\rU privitku se nalazi zatraženi transakcijski izvještaj."

	if err := SendEmail(req.Email, subject, body, attachment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email: " + err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"success": true})
}
