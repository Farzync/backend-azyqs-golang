package validators

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// Regex untuk validasi email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Regex untuk validasi username (huruf, angka, titik, tanpa titik di awal/akhir)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*$`)

// Regex sederhana untuk validasi panjang password
var passwordRegex = regexp.MustCompile(`^.{8,}$`)

// ValidateEmail memastikan email dalam format yang benar
func ValidateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return errors.New("invalid_email_format")
	}
	return nil
}

// ValidateUsername memastikan username sesuai aturan
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username_too_short")
	}
	if len(username) > 32 {
		return errors.New("username_too_long")
	}
	if !usernameRegex.MatchString(username) {
		return errors.New("invalid_username_format")
	}
	if strings.HasPrefix(username, ".") || strings.HasSuffix(username, ".") {
		return errors.New("username_cannot_start_or_end_with_dot")
	}
	return nil
}

// ValidateName memastikan nama sesuai aturan
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return errors.New("name_too_short")
	}
	if len(name) > 32 {
		return errors.New("name_too_long")
	}
	return nil
}

// ValidatePassword memastikan password memenuhi syarat
func ValidatePassword(password string) error {
	if !passwordRegex.MatchString(password) {
		return errors.New("password_too_short")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password_must_include_upper_lower_digit_special")
	}

	return nil
}
