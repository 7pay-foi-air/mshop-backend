package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/mshop/account-service/models"
)

type RegistrationRepository struct {
	db *sql.DB
}

func NewRegistrationRepository(db *sql.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) CreateUser(tx *sql.Tx, userUUID uuid.UUID, user models.RegistrationRequest, passwordHash string, dob time.Time) error {
	role := "cashier"
	if user.IsAdmin {
		role = "admin"
	}
	_, err := tx.Exec(`
		INSERT INTO user_account (
			uuid_user, first_name, last_name, username, email, phone_number,
			date_of_birth, address, password_hash, role, uuid_organisation 
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		userUUID,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.PhoneNumber,
		dob,
		user.Address,
		passwordHash,
		role,
		user.OrganisationUUID,
	)
	return err
}

func (r *RegistrationRepository) SetRecoveryToken(
	userID uuid.UUID,
	recoveryHash string,
) error {
	_, err := r.db.Exec(`
		UPDATE user_account
		SET
			recovery_token_hash = $1,
			updated_at = now()
		WHERE uuid_user = $2
		  AND deleted_at IS NULL
	`, recoveryHash, userID)

	return err
}

func (r *RegistrationRepository) ChangePasswordWithRecovery(
	userID uuid.UUID,
	newPasswordHash string,
) error {
	_, err := r.db.Exec(`
		UPDATE user_account
		SET
			password_hash = $1,
			updated_at = now(),
			lockout_counter = 0,
			is_locked = fase
		WHERE uuid_user = $2
		  AND deleted_at IS NULL
	`, newPasswordHash, userID)

	return err
}
