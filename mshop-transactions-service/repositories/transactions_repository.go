package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mshop/transactions-service/models"
)

type TransactionsRepository interface {
	CreateTransaction(tx *sql.Tx, req models.CreateTransactionRequest, user uuid.UUID, org uuid.UUID) (uuid.UUID, error)
	InsertTransactionItem(tx *sql.Tx, item models.TransactionItemRequest, txID uuid.UUID) (float64, error)
	FinalizeTransaction(tx *sql.Tx, txID uuid.UUID, total float64) error
	GetUserTransactions(userID uuid.UUID, orgID uuid.UUID, isAdmin bool, startDate, endDate string) (models.TransactionHistoryResponse, error)
}

type transactionsRepository struct {
	db *sql.DB
}

func NewTransactionsRepository(db *sql.DB) TransactionsRepository {
	return &transactionsRepository{db: db}
}

func (r *transactionsRepository) CreateTransaction(
	tx *sql.Tx,
	req models.CreateTransactionRequest,
	user uuid.UUID,
	org uuid.UUID,
) (uuid.UUID, error) {

	var id uuid.UUID

	err := tx.QueryRow(`
		INSERT INTO transaction (
			total_amount,
			currency,
			payment_method,
			is_successful,
			transaction_type,
			description,
			uuid_organisation,
			uuid_user
		)
		VALUES (0, $1, $2, false, 'Purchase', $3, $4, $5)
		RETURNING uuid_transaction
	`,
		req.Currency,
		req.PaymentMethod,
		req.Description,
		org,
		user,
	).Scan(&id)

	return id, err
}

func (r *transactionsRepository) InsertTransactionItem(
	tx *sql.Tx,
	item models.TransactionItemRequest,
	txID uuid.UUID,
) (float64, error) {

	if item.Quantity <= 0 || item.ItemPrice < 0 {
		return 0, errors.New("invalid item data")
	}

	subtotal := item.ItemPrice * float64(item.Quantity)

	_, err := tx.Exec(`
		INSERT INTO transaction_item (
			item_name,
			item_price,
			quantity,
			subtotal,
			uuid_item,
			uuid_transaction
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		item.ItemName,
		item.ItemPrice,
		item.Quantity,
		subtotal,
		item.UUIDItem,
		txID,
	)

	return subtotal, err
}

func (r *transactionsRepository) FinalizeTransaction(
	tx *sql.Tx,
	txID uuid.UUID,
	total float64,
) error {

	_, err := tx.Exec(`
		UPDATE transaction
		SET total_amount = $1,
		    is_successful = true,
		    completed_at = NOW()
		WHERE uuid_transaction = $2
	`, total, txID)

	return err
}

func (r *transactionsRepository) GetUserTransactions(userID uuid.UUID, orgID uuid.UUID, isAdmin bool, startDate, endDate string) (models.TransactionHistoryResponse, error) {
	var resp models.TransactionHistoryResponse

	query, args := r.buildTransactionQuery(userID, orgID, isAdmin, startDate, endDate)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return resp, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	transactions, err := r.scanTransactions(rows)
	if err != nil {
		return resp, err
	}

	resp.SuccessfulTransactions, resp.RefundedTransactions = r.categorizeTransactions(transactions)

	return resp, nil
}

func (r *transactionsRepository) buildTransactionQuery(userID uuid.UUID, orgID uuid.UUID, isAdmin bool, startDate, endDate string) (string, []interface{}) {
	query := `
		SELECT uuid_transaction, total_amount, currency, created_at, uuid_refund_to_transaction
		FROM transaction
		WHERE uuid_organisation = $1 AND is_successful = true`

	args := []interface{}{orgID}
	argPos := 2

	if !isAdmin {
		query += fmt.Sprintf(" AND uuid_user = $%d", argPos)
		args = append(args, userID)
		argPos++
	}

	if startDate != "" {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, startDate)
		argPos++
	}

	if endDate != "" {
		query += fmt.Sprintf(" AND created_at < $%d::date + interval '1 day'", argPos)
		args = append(args, endDate)
		argPos++
	}

	query += " ORDER BY created_at DESC"

	return query, args
}

func (r *transactionsRepository) scanTransactions(rows *sql.Rows) ([]models.TransactionHistory, error) {
	var transactions []models.TransactionHistory

	for rows.Next() {
		var t models.TransactionHistory
		err := rows.Scan(
			&t.UUIDTransaction,
			&t.TotalAmount,
			&t.Currency,
			&t.TransactionDate,
			&t.TransactionRefundID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

func (r *transactionsRepository) categorizeTransactions(transactions []models.TransactionHistory) ([]models.TransactionHistory, []models.TransactionHistory) {
	var successful, refunded []models.TransactionHistory

	for _, t := range transactions {
		if t.TransactionRefundID == nil {
			successful = append(successful, t)
		} else {
			refunded = append(refunded, t)
		}
	}

	return successful, refunded
}
