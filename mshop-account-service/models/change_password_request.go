package models

type ChangePasswordRequest struct {
	RecoveryToken string `json:"recovery_token" binding:"required"`
	NewPassword   string `json:"new_password" binding:"required,min=10"`
}
