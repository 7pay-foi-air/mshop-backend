package models

// swagger:model UserRegisterRequest
type UserRegisterRequest struct {
	FirstName   string `json:"first_name" example:"Ivan" validate:"required,min=3"`
	LastName    string `json:"last_name" example:"Ivić" validate:"required,min=3"`
	Username    string `json:"username" example:"ivan.ivic" validate:"required,min=6"`
	Email       string `json:"email" example:"ivan@example.com" validate:"required,email"`
	PhoneNumber string `json:"phone_number" example:"+38599111222" validate:"required"`
	Address     string `json:"address" example:"Savska cesta 14, Zagreb" validate:"required,min=10"`
	DateOfBirth string `json:"date_of_birth" example:"1990-05-20" validate:"required,datetime=2006-01-02"`
	IsAdmin     bool   `json:"is_admin" example:"true"`
}

// swagger:model OrganizationRegisterRequest
type OrganizationRegisterRequest struct {
	Name        string `json:"name" example:"mShop d.o.o." validate:"required,min=3"`
	OIB         string `json:"oib" example:"12345678901" validate:"required,oib"`
	Address     string `json:"address" example:"Savska cesta 123, Zagreb" validate:"required,min=10"`
	PhoneNumber string `json:"phone_number" example:"+38515555555" validate:"required,numeric"`
	Email       string `json:"email" example:"info@mshop.hr" validate:"required,email"`
}

// RegistrationRequest combines user and organization registration data.
// swagger:model RegistrationRequest
type RegistrationRequest struct {
	User         UserRegisterRequest         `json:"user"`
	Organization OrganizationRegisterRequest `json:"organization"`
}
