package repositories

import (
	"database/sql"

	"github.com/google/uuid"
)

func DeleteItemRepository(db *sql.DB, id uuid.UUID) (int64, error) {
	query := `
		UPDATE item
		SET deleted_at = NOW()
		WHERE uuid_item = $1
		AND deleted_at IS NULL
	`

	res, err := db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
