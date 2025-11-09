package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mshop/account-service/models"
)

var validate = NewValidator()

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("oib", validateOIB)
	return v
}

func ValidateRegistration(req models.RegistrationRequest) error {
	if err := ValidateUser(req.User); err != nil {
		return err
	}
	if err := ValidateOrganization(req.Organization); err != nil {
		return err
	}
	return nil
}

func ValidateUser(user models.UserRegisterRequest) error {
	if err := validate.Struct(user); err != nil {
		var messages []string
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, ValidateFormat(err).Error())
		}

		return fmt.Errorf(strings.Join(messages, ", "))
	}

	return nil
}

func ValidateOrganization(organization models.OrganizationRegisterRequest) error {
	if err := validate.Struct(organization); err != nil {
		var messages []string
		for _, err := range err.(validator.ValidationErrors) {
			messages = append(messages, ValidateFormat(err).Error())
		}

		return fmt.Errorf(strings.Join(messages, "; "))
	}

	return nil
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
		return fmt.Errorf("%s is not a valid datetime", err.Field())
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
