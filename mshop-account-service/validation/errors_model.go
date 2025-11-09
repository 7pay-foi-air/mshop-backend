package validation

type ErrorResponse struct {
	Code   int          `json:"code"`
	Status string       `json:"status"`
	Errors []FieldError `json:"errors"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}
