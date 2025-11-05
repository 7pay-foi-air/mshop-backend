package models

// swagger:model UserRegisterRequest
type UserRegisterRequest struct {
	FirstName   string `json:"first_name" example:"Ivan"`
	LastName    string `json:"last_name" example:"Ivić"`
	Username    string `json:"username" example:"ivan.ivic"`
	Email       string `json:"email" example:"ivan@example.com"`
	PhoneNumber string `json:"phone_number" example:"+38599111222"`
	Address     string `json:"address" example:"Savska cesta 14, Zagreb"`
	DateOfBirth string `json:"date_of_birth" example:"1990-05-20"`
	IsAdmin     bool   `json:"is_admin" example:"true"`
}

// swagger:model OrganizationRegisterRequest
type OrganizationRegisterRequest struct {
	Name        string `json:"name" example:"mShop d.o.o."`
	OIB         string `json:"oib" example:"12345678901"`
	Address     string `json:"address" example:"Savska cesta 123, Zagreb"`
	PhoneNumber string `json:"phone_number" example:"+38515555555"`
	Email       string `json:"email" example:"info@mshop.hr"`
}

// RegistrationRequest combines user and organization registration data.
// swagger:model RegistrationRequest
type RegistrationRequest struct {
	User         UserRegisterRequest         `json:"user"`
	Organization OrganizationRegisterRequest `json:"organization"`
}
