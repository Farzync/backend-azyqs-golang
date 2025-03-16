package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET_KEY = []byte("your_secret_key") // Replace with a secure secret key

// Custom error codes for JWT
var (
	ErrTokenExpired    = errors.New("token_expired")
	ErrTokenInvalid    = errors.New("token_invalid")
	ErrTokenMalformed  = errors.New("token_malformed")
	ErrTokenUnexpected = errors.New("token_unexpected_signing_method")
	ErrTokenPayload    = errors.New("invalid_token_payload")
)

// GenerateJWT generates a JWT token based on userID
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT validates the token and returns the userID if valid
func ValidateJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenUnexpected
		}
		return SECRET_KEY, nil
	})

	// Handle parsing errors with detailed responses
	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			switch {
			case validationErr.Errors&jwt.ValidationErrorMalformed != 0:
				return 0, ErrTokenMalformed
			case validationErr.Errors&jwt.ValidationErrorExpired != 0:
				return 0, ErrTokenExpired
			default:
				return 0, ErrTokenInvalid
			}
		}
		return 0, ErrTokenInvalid
	}

	// Validate claims and extract userID
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return uint(userIDFloat), nil
		}
		return 0, ErrTokenPayload
	}

	return 0, ErrTokenInvalid
}
