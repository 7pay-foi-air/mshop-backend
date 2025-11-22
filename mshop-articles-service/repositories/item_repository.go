package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mshop/articles-service/models"
)

type ItemRepository interface {
	GetItems(ids []uuid.UUID) ([]models.ItemResponse, error)
	CreateItem(item models.Item) (uuid.UUID, error)
	DeleteItem(id uuid.UUID) (int64, error)
}

type itemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

func scanItem(rows *sql.Rows) (models.ItemResponse, error) {
	var item models.ItemResponse

	err := rows.Scan(
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
	)

	return item, err
}

func (r *itemRepository) GetItems(ids []uuid.UUID) ([]models.ItemResponse, error) {
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
		FROM items
		WHERE deleted_at IS NULL
	`

	var (
		rows *sql.Rows
		err  error
	)

	if len(ids) > 0 {
		query := baseQuery + " AND uuid_item = ANY($1)"
		rows, err = r.db.Query(query, pq.Array(ids))
	} else {
		rows, err = r.db.Query(baseQuery)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	defer rows.Close()

	var items []models.ItemResponse

	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return items, nil
}

func (r *itemRepository) DeleteItem(id uuid.UUID) (int64, error) {
	query := `
		UPDATE items
		SET deleted_at = NOW()
		WHERE uuid_item = $1
		  AND deleted_at IS NULL
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete item: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return affected, nil
}

func (r *itemRepository) CreateItem(item models.Item) (uuid.UUID, error) {
	query := `
		INSERT INTO item (
			uuid_organisation,
			name,
			description,
			price,
			currency,
			sku,
			stock_quantity,
			image_url
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		RETURNING uuid_item
	`

	var newID uuid.UUID
	err := r.db.QueryRow(query,
		item.UUIDOrganisation,
		item.Name,
		item.Description,
		item.Price,
		item.Currency,
		item.SKU,
		item.StockQuantity,
		item.ImageURL,
	).Scan(&newID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert item: %w", err)
	}

	return newID, nil
}
