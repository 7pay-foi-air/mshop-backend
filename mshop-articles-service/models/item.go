package models

import "github.com/google/uuid"

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
