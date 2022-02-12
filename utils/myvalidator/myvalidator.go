package myvalidator

import (
	"unicode"
	"unicode/utf8"

	"gopkg.in/go-playground/validator.v9"
)

// ValidateMyVal implements validator.Func
func ValidateUserName(fl validator.FieldLevel) bool {

	username := fl.Field().String()

	if utf8.RuneCountInString(username) > 30 || utf8.RuneCountInString(username) < 3 {
		return false
	}

	for _, c := range username {
		switch {
		case unicode.IsNumber(c):
			break
		case c == '-' || c == '_' || c == '.' || c == '@':
			break
		case unicode.IsLetter(c):
			break
		default:
			return false
		}
	}

	return true
}
