package models

// LoginRequest represents a simple login payload
// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" example:"iivanic7" validate:"required,min=6"`
	Password string `json:"password" example:"password123" validate:"required,min=10"`
}
