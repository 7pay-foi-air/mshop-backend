package repositories

import (
	"database/sql"

	"github.com/mshop/articles-service/models"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// Create inserts a new item into the database
func (r *ItemRepository) CreateItem(item *models.ItemCreateRequest) (string, error) {
	var id string
	query := `
		INSERT INTO item (
			uuid_organisation, name, description, price, currency, sku, stock_quantity, image_url
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING uuid_item
	`
	err := r.db.QueryRow(
		query,
		item.OrganisationID,
		item.Name,
		item.Description,
		item.Price,
		item.Currency,
		item.SKU,
		item.StockQuantity,
		item.ImageURL,
	).Scan(&id)

	return id, err
}
