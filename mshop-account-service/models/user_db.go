package models

import (
	"time"

	"github.com/google/uuid"
)

type UserDB struct {
	UUID                       uuid.UUID  `json:"uuid_user"`
	FirstName                  string     `json:"first_name"`
	LastName                   string     `json:"last_name"`
	Username                   string     `json:"username"`
	Email                      string     `json:"email"`
	IsEmailVerified            bool       `json:"is_email_verified"`
	PhoneNumber                string     `json:"phone_number"`
	DateOfBirth                time.Time  `json:"date_of_birth"`
	Address                    string     `json:"address"`
	PasswordHash               string     `json:"password_hash"`
	LastLoginAt                *time.Time `json:"last_login_at"`
	RecoveryTokenHash          *string    `json:"recovery_token_hash"`
	CreatedAt                  time.Time  `json:"created_at"`
	UpdatedAt                  time.Time  `json:"updated_at"`
	DeletedAt                  *time.Time `json:"deleted_at"`
	IsActive                   bool       `json:"is_active"`
	EmailVerificationTokenHash *string    `json:"email_verification_token_hash"`
	EmailVerificationExpiresAt *time.Time `json:"email_verification_expires_at"`
	Role                       string     `json:"role"`
	OrganisationUUID           *uuid.UUID `json:"uuid_organisation"`
}
