package repositories

import (
	"database/sql"
)

type SecurityRepository struct {
	db *sql.DB
}

func NewSecurityRepository(db *sql.DB) *SecurityRepository {
	return &SecurityRepository{db: db}
}

func (r *SecurityRepository) SetSecurityQuestions(userID, questionsHash, recoveryLocation string) error {
	query := `
		UPDATE user_account
		SET security_questions_hash = $1,
		    recovery_code_location = $2,
		    updated_at = NOW()
		WHERE uuid_user = $3 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, questionsHash, recoveryLocation, userID)
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

func (r *SecurityRepository) GetSecurityQuestions(userID string) (string, error) {
	var questionsHash string

	query := `
		SELECT security_questions_hash
		FROM user_account
		WHERE uuid_user = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, userID).Scan(&questionsHash)
	if err != nil {
		return "", err
	}

	return questionsHash, nil
}
