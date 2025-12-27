package validation

import "unicode"

func ValidatePassword(password string) bool {
	if len(password) < 10 {
		return false
	}

	hasLetter := false
	hasNumber := false

	for _, c := range password {
		if unicode.IsLetter(c) {
			hasLetter = true
		}
		if unicode.IsNumber(c) {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}
