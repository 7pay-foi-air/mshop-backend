package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mshop/articles-service/models"
)

type ItemRepository struct {
	db *sql.DB
}

func GetItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetItems(uuids []uuid.UUID) ([]models.ItemResponse, error) {
	items := []models.ItemResponse{}

	baseQuery := `
		SELECT 
			uuid_item,
			uuid_organisation,
			name,
			description,
			price,
			currency,
			sku,
			stock_quantity,
			image_url,
			is_active,
			created_at,
			updated_at
		FROM item
		WHERE deleted_at IS NULL
	`

	var (
		rows *sql.Rows
		err  error
	)

	if len(uuids) > 0 {
		query := baseQuery + " AND uuid_item = ANY($1)"
		rows, err = r.db.Query(query, pq.Array(uuids))
	} else {
		rows, err = r.db.Query(baseQuery)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ItemResponse

		if err := rows.Scan(
			&item.UUIDItem,
			&item.UUIDOrganisation,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Currency,
			&item.SKU,
			&item.StockQuantity,
			&item.ImageURL,
			&item.IsActive,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
