package repositories

import (
	"database/sql"

	"github.com/google/uuid"
)

type SecurityRepository struct {
	db *sql.DB
}

func NewSecurityRepository(db *sql.DB) *SecurityRepository {
	return &SecurityRepository{db: db}
}

func (r *SecurityRepository) SetSecurityQuestions(username, questionsHash, recoveryLocation string) error {
	query := `
		UPDATE user_account
		SET security_questions_hash = $1,
		    recovery_code_location = $2,
		    updated_at = NOW()
		WHERE username = $3 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, questionsHash, recoveryLocation, username)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SecurityRepository) GetSecurityQuestions(username string) (string, string, error) {
	var questionsHash string
	var recoveryCodeLocation string

	query := `
		SELECT security_questions_hash, recovery_code_location
		FROM user_account
		WHERE username = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, username).
		Scan(&questionsHash, &recoveryCodeLocation)
	if err != nil {
		return "", "", err
	}

	return questionsHash, recoveryCodeLocation, nil
}

func (r *SecurityRepository) GetRecoveryLocationByUserUUID(userUUID uuid.UUID) (string, bool, bool, error) {
	var location sql.NullString
	var isActive bool
	var isLocked bool

	query := `
		SELECT recovery_code_location, is_active, is_locked
		FROM user_account
		WHERE uuid_user = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, userUUID).Scan(&location, &isActive, &isLocked)
	if err != nil {
		// ako nema reda, Scan vraća sql.ErrNoRows
		return "", false, false, err
	}

	if !location.Valid {
		return "", isActive, isLocked, nil
	}

	return location.String, isActive, isLocked, nil
}
