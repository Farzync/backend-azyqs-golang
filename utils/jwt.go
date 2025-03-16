package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET_KEY = []byte("your_secret_key") // Ganti dengan secret key yang aman

// GenerateJWT membuat token JWT berdasarkan userID
func GenerateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT memvalidasi token dan mengembalikan userID jika valid
func ValidateJWT(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan signing method yang digunakan adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return SECRET_KEY, nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Perlu konversi dari float64 ke uint
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return uint(userIDFloat), nil
		}
		return 0, errors.New("Invalid token payload")
	}
	return 0, errors.New("Token tidak valid")
}
