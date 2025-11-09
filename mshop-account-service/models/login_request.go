package models

// LoginRequest represents a simple login payload
// swagger:model LoginRequest
type LoginRequest struct {
	Username string `json:"username" example:"ivan.ivic"`
	Password string `json:"password" example:"test123"`
}
