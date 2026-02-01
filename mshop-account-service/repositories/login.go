package repositories

import (
	"database/sql"
	"errors"

	"github.com/mshop/account-service/models"
)

type LoginRepository struct {
	db *sql.DB
}

func NewLoginRepository(db *sql.DB) *LoginRepository {
	return &LoginRepository{db: db}
}

func (r *LoginRepository) GetUserByUsername(username string) (*models.UserDB, error) {
	user := models.UserDB{}

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
        password_hash,
        last_login_at,
        recovery_token_hash,
        created_at,
        updated_at,
        deleted_at,
        is_active,
        role,
        uuid_organisation,
		recovery_code_location,
		security_questions_hash,
		is_locked,
		lockout_counter
	FROM user_account
	WHERE username = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRow(query, username)
	err := row.Scan(
		&user.UUID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.DateOfBirth,
		&user.Address,
		&user.PasswordHash,
		&user.LastLoginAt,
		&user.RecoveryTokenHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.IsActive,
		&user.Role,
		&user.OrganisationUUID,
		&user.RecoveryCodeLocation,
		&user.SecurityQuestionsHash,
		&user.IsLocked,
		&user.LockoutCounter,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *LoginRepository) GetUserByID(userID string) (*models.UserDB, error) {
	user := models.UserDB{}

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
        password_hash,
        last_login_at,
        recovery_token_hash,
        created_at,
        updated_at,
        deleted_at,
        is_active,
        role,
        uuid_organisation
	FROM user_account
	WHERE uuid_user = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRow(query, userID)
	err := row.Scan(
		&user.UUID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.DateOfBirth,
		&user.Address,
		&user.PasswordHash,
		&user.LastLoginAt,
		&user.RecoveryTokenHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.IsActive,
		&user.Role,
		&user.OrganisationUUID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
func (r *LoginRepository) UpdateLoginSecurity(user *models.UserDB) error {
	query := `
	UPDATE user_account
	SET 
		lockout_counter = $1,
		is_locked = $2,
		updated_at = NOW()
	WHERE uuid_user = $3
	`

	_, err := r.db.Exec(
		query,
		user.LockoutCounter,
		user.IsLocked,
		user.UUID,
	)

	return err
}
