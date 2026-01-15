package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mshop/account-service/models"
)

type UserRepository interface {
	GetUsers(ids []uuid.UUID) ([]models.UserDB, error)
	GetUserByID(id uuid.UUID) (*models.UserDB, error)
	UpdateUser(userUUID uuid.UUID, updates map[string]interface{}) (*models.UserDB, error)
	UpdateUserByAdmin(userUUID uuid.UUID, updates map[string]interface{}) (*models.UserDB, error)
	DeleteUser(userUUID uuid.UUID) (int64, error)
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

func scanUserRow(row *sql.Row) (models.UserDB, error) {
	var user models.UserDB

	err := row.Scan(
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

func (r *userRepository) GetUserByID(id uuid.UUID) (*models.UserDB, error) {
	query := `
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
		WHERE uuid_user = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRow(query, id)
	user, err := scanUserRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(userUUID uuid.UUID, updates map[string]interface{}) (*models.UserDB, error) {
	_, err := r.GetUserByID(userUUID)
	if err != nil {
		return nil, err
	}

	allowedFields := map[string]bool{
		"first_name":    true,
		"last_name":     true,
		"date_of_birth": true,
		"phone_number":  true,
		"address":       true,
		"email":         true,
	}

	setClauses := []string{}
	args := []interface{}{}
	paramCount := 1

	for field, value := range updates {
		if !allowedFields[field] {
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, paramCount))
		args = append(args, value)
		paramCount++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramCount))
	args = append(args, time.Now())
	paramCount++

	args = append(args, userUUID)

	query := fmt.Sprintf(`
		UPDATE user_account
		SET %s
		WHERE uuid_user = $%d AND deleted_at IS NULL
	`, strings.Join(setClauses, ", "), paramCount)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("duplicate value: %s", pqErr.Constraint)
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found or already deleted")
	}

	return r.GetUserByID(userUUID)
}

func (r *userRepository) UpdateUserByAdmin(userUUID uuid.UUID, updates map[string]interface{}) (*models.UserDB, error) {
	_, err := r.GetUserByID(userUUID)
	if err != nil {
		return nil, err
	}

	allowedFields := map[string]bool{
		"first_name":    true,
		"last_name":     true,
		"date_of_birth": true,
		"phone_number":  true,
		"address":       true,
		"email":         true,
		"username":      true,
		"role":          true,
		"is_active":     true,
	}

	setClauses := []string{}
	args := []interface{}{}
	paramCount := 1

	for field, value := range updates {
		if !allowedFields[field] {
			continue
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, paramCount))
		args = append(args, value)
		paramCount++
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no valid fields to update")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", paramCount))
	args = append(args, time.Now())
	paramCount++

	args = append(args, userUUID)

	query := fmt.Sprintf(`
		UPDATE user_account
		SET %s
		WHERE uuid_user = $%d AND deleted_at IS NULL
	`, strings.Join(setClauses, ", "), paramCount)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return nil, fmt.Errorf("duplicate value: %s", pqErr.Constraint)
			}
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found or already deleted")
	}

	return r.GetUserByID(userUUID)
}

func (r *userRepository) DeleteUser(id uuid.UUID) (int64, error) {
	query := `
		UPDATE user_account
		SET deleted_at = NOW(),
    		is_active = FALSE
		WHERE uuid_user = $1
		AND deleted_at IS NULL
	`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get affected rows: %w", err)
	}

	return affected, nil
}
