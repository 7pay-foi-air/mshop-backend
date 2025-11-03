package models

type UserRegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	IsAdmin     bool   `json:"is_admin"`
}
