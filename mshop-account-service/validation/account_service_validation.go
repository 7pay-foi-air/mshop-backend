package validation

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/mshop/account-service/models"
)

var validate = NewValidator()

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("oib", validateOIB)
	return v
}

func ValidateRegistration(req models.RegistrationRequest) []string {
	var messages []string

	if err := ValidateUser(req.User); err != nil {
		messages = append(messages, ValidateUser(req.User)...)
	}
	if err := ValidateOrganization(req.Organization); err != nil {
		messages = append(messages, ValidateOrganization(req.Organization)...)
	}

	return messages
}

func ValidateUser(user models.UserRegisterRequest) []string {
	var messages []string

	if err := validate.Struct(user); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, ValidateFormat(err).Error())
		}
	}

	return messages
}

func ValidateOrganization(organization models.OrganizationRegisterRequest) []string {
	var messages []string

	if err := validate.Struct(organization); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, ValidateFormat(err).Error())
		}
	}

	return messages
}

func validateOIB(fl validator.FieldLevel) bool {
	oib := fl.Field().String()

	matched, _ := regexp.MatchString("^[0-9]{11}$", oib)
	return matched
}

func ValidateFormat(err validator.FieldError) error {
	switch err.Tag() {
	case "required":
		return fmt.Errorf("%s is required", err.Field())
	case "email":
		return fmt.Errorf("%s is not a valid email", err.Field())
	case "numeric":
		return fmt.Errorf("%s not a valid number", err.Field())
	case "datetime":
		return fmt.Errorf("%s is not a valid date", err.Field())
	case "min":
		return fmt.Errorf("%s is too short", err.Field())
	case "len":
		return fmt.Errorf("%s is not the proper length", err.Field())
	case "oib":
		return fmt.Errorf("%s must contain exactly 11 digits", err.Field())
	default:
		return fmt.Errorf("%s is invalid", err.Field())
	}
}
