package models

type ItemCreateRequest struct {
	Name           string  `json:"name" example:"Laptop"`
	Description    string  `json:"description" example:"Powerful gaming laptop"`
	Price          float64 `json:"price" example:"1299.99"`
	Currency       string  `json:"currency" example:"EUR"`
	SKU            string  `json:"sku" example:"ABC-123"`
	StockQuantity  int     `json:"stock_quantity" example:"10"`
	ImageURL       string  `json:"image_url" example:"https://example.com/image.png"`
	OrganisationID string  `json:"uuid_organisation" example:"550e8400-e29b-41d4-a716-446655440000"`
}
