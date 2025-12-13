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
	v.RegisterValidation("telephone", validateTelephone)
	return v
}

func ValidateRegistration(req models.RegistrationRequest) []FieldError {
	var errors []FieldError

	if err := ValidateModel(req); err != nil {
		errors = append(errors, ValidateModel(req)...)
	}
	return errors
}

func ValidateLogin(req models.LoginRequest) []FieldError {
	return ValidateModel(req)
}

func ValidateModel(model any) []FieldError {
	var fieldErrors []FieldError

	if err := validate.Struct(model); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldErrors = append(fieldErrors, ValidateFormat(err))
		}
	}

	return fieldErrors
}

func validateOIB(fl validator.FieldLevel) bool {
	oib := fl.Field().String()

	matched, _ := regexp.MatchString("^[0-9]{11}$", oib)
	return matched
}

func validateTelephone(fl validator.FieldLevel) bool {
	telephone := fl.Field().String()
	matched, _ := regexp.MatchString(`^\+?[0-9]{1,4}?[-.\s]?\(?[0-9]{1,3}?\)?[-.\s]?[0-9]{3,4}[-.\s]?[0-9]{3,4}$`, telephone)
	return matched
}

func ValidateFormat(err validator.FieldError) FieldError {
	tag := err.Tag()

	if template, exists := ValidationMessages[tag]; exists {
		return FieldError{
			Field:   err.Field(),
			Message: fmt.Sprintf(template.Message, err.Field()),
			Reason:  template.Reason,
		}
	}

	return FieldError{
		Field:   err.Field(),
		Message: fmt.Sprintf("%s failed validation rule '%s'", err.Field(), tag),
		Reason:  "The provided value failed an unrecognized validation rule",
	}
}
