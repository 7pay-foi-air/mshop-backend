package models

// UpdateUserRequest represents basic user info that regular users can update
type UpdateUserRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Email       *string `json:"email"`
}

// AdminUpdateUserRequest represents all fields that admins can update
type AdminUpdateUserRequest struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Email       *string `json:"email"`
	Username    *string `json:"username"`
	Role        *string `json:"role"`
	IsActive    *bool   `json:"is_active"`
}
