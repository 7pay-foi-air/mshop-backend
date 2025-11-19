package models

import "github.com/google/uuid"

// swagger:model RegistrationRequest
type RegistrationRequest struct {
	FirstName        string    `json:"first_name" example:"Ivan" validate:"required,min=3"`
	LastName         string    `json:"last_name" example:"Ivić" validate:"required,min=3"`
	Username         string    `json:"username" example:"ivan.ivic" validate:"required,min=6"`
	Email            string    `json:"email" example:"ivan@example.com" validate:"required,email"`
	PhoneNumber      string    `json:"phone_number" example:"+38599111222" validate:"required,telephone"`
	Address          string    `json:"address" example:"Savska cesta 14, Zagreb" validate:"required,min=10"`
	DateOfBirth      string    `json:"date_of_birth" example:"1990-05-20" validate:"required,datetime=2006-01-02"`
	IsAdmin          bool      `json:"is_admin" example:"true"`
	OrganisationUUID uuid.UUID `json:"organisation_uuid" example:"02f2c243-6c29-4f21-a98c-955372bc6297" validate:"required"`
}
