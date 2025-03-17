package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
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
func ValidateJWT(tokenString string) (uuid.UUID, error) {
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
				return uuid.Nil, ErrTokenMalformed
			case validationErr.Errors&jwt.ValidationErrorExpired != 0:
				return uuid.Nil, ErrTokenExpired
			default:
				return uuid.Nil, ErrTokenInvalid
			}
		}
		return uuid.Nil, ErrTokenInvalid
	}

	// Validate claims and extract userID
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDStr, ok := claims["user_id"].(string); ok {
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				return uuid.Nil, ErrTokenPayload
			}
			return userID, nil
		}
		return uuid.Nil, ErrTokenPayload
	}

	return uuid.Nil, ErrTokenInvalid
}