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
        is_email_verified,
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
        email_verification_token_hash,
        email_verification_expires_at,
        role,
        uuid_organisation
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
		&user.IsEmailVerified,
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
		&user.EmailVerificationTokenHash,
		&user.EmailVerificationExpiresAt,
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
