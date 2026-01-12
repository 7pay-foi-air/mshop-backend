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

	GetTransactionByID(txID uuid.UUID) (models.TransactionHistory, error)
	CreateRefundTransaction(tx *sql.Tx, original models.TransactionHistory, user uuid.UUID, description string) (uuid.UUID, error)
	GetTransactionItems(txID uuid.UUID) ([]models.TransactionItemDetail, error)
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

func (r *transactionsRepository) GetTransactionByID(txID uuid.UUID) (models.TransactionHistory, error) {
	var t models.TransactionHistory

	err := r.db.QueryRow(`
		SELECT uuid_transaction, total_amount, currency, created_at,
		       uuid_refund_to_transaction, payment_method, uuid_organisation,
		       transaction_type, uuid_user
		FROM transaction
		WHERE uuid_transaction = $1 AND is_successful = true
	`, txID).Scan(
		&t.UUIDTransaction,
		&t.TotalAmount,
		&t.Currency,
		&t.TransactionDate,
		&t.TransactionRefundID,
		&t.PaymentMethod,
		&t.UUIDOrganisation,
		&t.TransactionType,
		&t.UUIDUser,
	)

	return t, err
}

func (r *transactionsRepository) CreateRefundTransaction(
	tx *sql.Tx,
	original models.TransactionHistory,
	user uuid.UUID,
	description string,
) (uuid.UUID, error) {

	var refundID uuid.UUID

	err := tx.QueryRow(`
	INSERT INTO transaction (
		total_amount,
		currency,
		payment_method,
		is_successful,
		transaction_type,
		description,
		uuid_refund_to_transaction,
		uuid_user,
		uuid_organisation,
		completed_at
	)
	VALUES ($1, $2, $3, true, 'Refund', $4, $5, $6, $7, NOW()
	)
	RETURNING uuid_transaction
`,
		original.TotalAmount,
		original.Currency,
		original.PaymentMethod,
		description,
		original.UUIDTransaction,
		user,
		original.UUIDOrganisation,
	).Scan(&refundID)

	return refundID, err
}

func (r *transactionsRepository) GetTransactionItems(txID uuid.UUID) ([]models.TransactionItemDetail, error) {
	rows, err := r.db.Query(`
		SELECT uuid_item, item_name, item_price, quantity, subtotal
		FROM transaction_item
		WHERE uuid_transaction = $1
		ORDER BY item_name
	`, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transaction items: %w", err)
	}
	defer rows.Close()

	items := []models.TransactionItemDetail{}

	for rows.Next() {
		var it models.TransactionItemDetail
		if err := rows.Scan(&it.UUIDItem, &it.ItemName, &it.ItemPrice, &it.Quantity, &it.Subtotal); err != nil {
			return nil, fmt.Errorf("failed to scan transaction item: %w", err)
		}
		items = append(items, it)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return items, nil
}
