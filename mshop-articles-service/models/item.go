package models

import (
	"time"

	"github.com/google/uuid"
)

// swagger:model ItemResponse
type ItemResponse struct {
	UUIDItem         uuid.UUID `json:"uuid_item"`
	UUIDOrganisation uuid.UUID `json:"uuid_organisation"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`
	SKU              *string   `json:"sku"`
	StockQuantity    int       `json:"stock_quantity"`
	ImageURL         *string   `json:"image_url"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}

type ItemCreateRequest struct {
	ItemName      string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Currency      Currency `json:"currency"`
	Sku           *string  `json:"sku"`
	StockQuantity uint64   `json:"stock_quantity"`
}

type ItemUpdateRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	Currency      Currency `json:"currency"`
	SKU           *string  `json:"sku"`
	StockQuantity uint64   `json:"stock_quantity"`
}

type Item struct {
	UUIDItem         uuid.UUID
	UUIDOrganisation uuid.UUID
	Name             string
	Description      string
	Price            float64
	Currency         string
	SKU              *string
	StockQuantity    int
	ImageURL         *string
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
