package models

// LoginRequest represents a simple login payload
// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" example:"ivan.ivic" validate:"required,min=6"`
	Password string `json:"password" example:"test123456" validate:"required,min=10"`
}
