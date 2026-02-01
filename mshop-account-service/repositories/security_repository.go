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

func (r *SecurityRepository) GetSecurityQuestions(username string) (string, error) {
	var questionsHash string

	query := `
		SELECT security_questions_hash
		FROM user_account
		WHERE username = $1 AND deleted_at IS NULL
	`

	err := r.db.QueryRow(query, username).Scan(&questionsHash)
	if err != nil {
		return "", err
	}

	return questionsHash, nil
}
