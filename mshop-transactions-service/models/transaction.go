package models

import (
	"github.com/google/uuid"
)

// swagger:model TransactionItemRequest
type TransactionItemRequest struct {
	UUIDItem  uuid.UUID `json:"uuid_item" example:"7c0d9f7a-b48e-4bb9-a8f3-4946b278c321"`
	ItemName  string    `json:"item_name" example:"Premium USB-C Kabel 1m"`
	ItemPrice float64   `json:"item_price" example:"9.99"`
	Quantity  int       `json:"quantity" example:"1"`
}

// swagger:model CreateTransactionRequest
type CreateTransactionRequest struct {
	PaymentMethod string                   `json:"payment_method" example:"credit_card"`
	Currency      string                   `json:"currency" example:"EUR"`
	Description   string                   `json:"description" example:"Kupnja USB-C kabela"`
	Items         []TransactionItemRequest `json:"items"`
}

// swagger:model TransactionResponse
type TransactionResponse struct {
	UUIDTransaction uuid.UUID `json:"uuid_transaction"`
	TotalAmount     float64   `json:"total_amount"`
	Currency        string    `json:"currency"`
	IsSuccessful    bool      `json:"is_successful"`
}

// swagger:model TransactionHistory
type TransactionHistory struct {
	UUIDTransaction     uuid.UUID  `json:"uuid_transaction"`
	TotalAmount         float64    `json:"total_amount"`
	Currency            string     `json:"currency"`
	TransactionDate     string     `json:"transaction_date"`
	TransactionRefundID *uuid.UUID `json:"transaction_refund_id"`
}

// swagger:model TransactionHistoryResponse
type TransactionHistoryResponse struct {
	SuccessfulTransactions []TransactionHistory `json:"successful_transactions"`
	RefundedTransactions   []TransactionHistory `json:"refunded_transactions"`
}
