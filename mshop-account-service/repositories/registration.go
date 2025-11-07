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

func (r *RegistrationRepository) CreateOrganisation(tx *sql.Tx, orgUUID uuid.UUID, org models.OrganizationRegisterRequest) error {
	_, err := tx.Exec(`
		INSERT INTO organisation (uuid_organisation, name, oib, address, contact_email, contact_phone)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		orgUUID,
		org.Name,
		org.OIB,
		org.Address,
		org.Email,
		org.PhoneNumber,
	)
	return err
}

func (r *RegistrationRepository) CreateUser(tx *sql.Tx, userUUID, orgUUID uuid.UUID, user models.UserRegisterRequest, passwordHash string, dob time.Time) error {
	_, err := tx.Exec(`
		INSERT INTO user_account (
			uuid_user, first_name, last_name, username, email, phone_number,
			date_of_birth, address, password_hash, uuid_organisation
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		orgUUID,
	)
	return err
}
