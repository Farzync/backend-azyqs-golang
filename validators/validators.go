package validators

import (
	"errors"
	"regexp"
	"strings"
)

// Regex untuk validasi email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Regex untuk validasi username (huruf, angka, titik, tanpa titik di awal/akhir)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)*$`)

// Regex untuk validasi password (minimal 8 karakter, harus ada huruf besar, kecil, angka, dan simbol)
var passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$`)

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

// ValidatePassword memastikan password memenuhi syarat
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password_too_short")
	}
	if !passwordRegex.MatchString(password) {
		return errors.New("password_must_include_upper_lower_digit_special")
	}
	return nil
}
