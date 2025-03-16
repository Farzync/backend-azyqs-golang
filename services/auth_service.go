package services

import (
	"azyqs-auth-systems/config"
	serviceErrors "azyqs-auth-systems/errors"
	"azyqs-auth-systems/models"
	"azyqs-auth-systems/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// RegisterUser registers a new user
func RegisterUser(username, name, email, password string) error {
	var count int64
	config.DB.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return serviceErrors.ErrDuplicateRecord
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return serviceErrors.ErrPasswordHash
	}

	user := models.User{
		Username: username,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return serviceErrors.ErrDuplicateRecord
		}
		return err
	}
	return nil
}

// LoginUser authenticates a user and returns a JWT token
func LoginUser(username, password string) (string, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", serviceErrors.ErrUserNotFound
		}
		return "", err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", serviceErrors.ErrInvalidPassword
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
