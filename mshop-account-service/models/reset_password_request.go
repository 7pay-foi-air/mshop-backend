package models

type ResetPasswordRequest struct {
	Username      string `json:"username" binding:"required"`
	RecoveryToken string `json:"recovery_token" binding:"required"`
	NewPassword   string `json:"new_password" binding:"required,min=10"`
}
