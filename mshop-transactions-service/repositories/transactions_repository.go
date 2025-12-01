package repositories

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/mshop/transactions-service/models"
)

type TransactionsRepository interface {
	CreateTransaction(tx *sql.Tx, req models.CreateTransactionRequest, user uuid.UUID, org uuid.UUID) (uuid.UUID, error)
	InsertTransactionItem(tx *sql.Tx, item models.TransactionItemRequest, txID uuid.UUID) (float64, error)
	FinalizeTransaction(tx *sql.Tx, txID uuid.UUID, total float64) error
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
