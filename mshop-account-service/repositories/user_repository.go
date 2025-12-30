package repositories

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mshop/account-service/models"
)

type UserRepository interface {
	GetUsers(ids []uuid.UUID) ([]models.UserDB, error)
	//UpdateUser()
	//DeleteUser()
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func scanUser(rows *sql.Rows) (models.UserDB, error) {
	var user models.UserDB

	err := rows.Scan(
		&user.UUID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.DateOfBirth,
		&user.Address,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return user, err
}

func (r *userRepository) GetUsers(ids []uuid.UUID) ([]models.UserDB, error) {
	baseQuery := `
		SELECT 
			uuid_user,
			first_name,
			last_name,
			username,
			email,
			phone_number,
			date_of_birth,
			address,
			role,
			is_active,
			created_at,
			updated_at
		FROM user_account
		WHERE deleted_at IS NULL
	`

	var (
		rows *sql.Rows
		err  error
	)

	if len(ids) > 0 {
		query := baseQuery + " AND uuid_user = ANY($1)"
		rows, err = r.db.Query(query, pq.Array(ids))
	} else {
		rows, err = r.db.Query(baseQuery)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []models.UserDB

	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}
