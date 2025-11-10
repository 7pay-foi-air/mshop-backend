package validation

var ValidationMessages = map[string]FieldError{
	"required": {
		Message: "%s is required",
		Reason:  "This field is mandatory but missing from the request",
	},
	"email": {
		Message: "%s is not a valid email",
		Reason:  "The provided value does not match standard email format",
	},
	"numeric": {
		Message: "%s must be a number",
		Reason:  "Only numeric characters are allowed for this field",
	},
	"datetime": {
		Message: "%s is not a valid date",
		Reason:  "The date format must be YYYY-MM-DD",
	},
	"min": {
		Message: "%s is too short",
		Reason:  "This field does not meet minimum length requirements",
	},
	"len": {
		Message: "%s has invalid length",
		Reason:  "Expected exact length according to validation rule",
	},
	"oib": {
		Message: "%s must contain exactly 11 digits",
		Reason:  "OIB requires exactly 11 numeric digits",
	},
	"telephone": {
		Message: "%s is not a valid telephone number",
		Reason:  "The provided value does not match a standard telephone number",
	},
}
