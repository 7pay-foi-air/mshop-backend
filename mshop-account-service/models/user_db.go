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
	Role                       string     `json:"role"`
	OrganisationUUID           *uuid.UUID `json:"uuid_organisation"`
}
