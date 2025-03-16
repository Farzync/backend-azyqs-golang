package services

import (
	"azyqs-auth-systems/config"
	"azyqs-auth-systems/models"
	"azyqs-auth-systems/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// Custom error codes
var (
	ErrUsernameTaken     = errors.New("username_already_taken")
	ErrEmailTaken        = errors.New("email_already_taken")
	ErrUserNotFound      = errors.New("user_not_found")
	ErrInvalidPassword   = errors.New("invalid_password")
	ErrPasswordMismatch  = errors.New("password_mismatch")
	ErrPasswordHash      = errors.New("password_hash_error")
	ErrUserDeleteFailed  = errors.New("user_delete_failed")
	ErrUserUpdateFailed  = errors.New("user_update_failed")
	ErrInvalidInput      = errors.New("invalid_input")
	ErrDuplicateRecord   = errors.New("duplicate_record")
)

// RegisterUser registers a new user
func RegisterUser(username, name, email, password string) error {
	// Check if username or email already exists
	var count int64
	config.DB.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return ErrDuplicateRecord
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return ErrPasswordHash
	}

	user := models.User{
		Username: username,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		// Handle duplicate record error
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return ErrDuplicateRecord
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
			return "", ErrUserNotFound
		}
		return "", err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", ErrInvalidPassword
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetUserByID fetches user data by ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile updates username, name, and email for a user
func UpdateUserProfile(userID uint, newUsername, newName, newEmail string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}

	// Check for username uniqueness if changed
	if newUsername != user.Username {
		var count int64
		config.DB.Model(&models.User{}).
			Where("username = ? AND id != ?", newUsername, userID).
			Count(&count)
		if count > 0 {
			return ErrUsernameTaken
		}
		user.Username = newUsername
	}

	// Check for email uniqueness if changed
	if newEmail != user.Email {
		var count int64
		config.DB.Model(&models.User{}).
			Where("email = ? AND id != ?", newEmail, userID).
			Count(&count)
		if count > 0 {
			return ErrEmailTaken
		}
		user.Email = newEmail
	}

	user.Name = newName

	if err := config.DB.Save(&user).Error; err != nil {
		return ErrUserUpdateFailed
	}
	return nil
}

// DeleteUser deletes a user after password confirmation
func DeleteUser(userID uint, password string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return ErrPasswordMismatch
	}
	if err := config.DB.Delete(&user).Error; err != nil {
		return ErrUserDeleteFailed
	}
	return nil
}

// ChangeUserPassword changes a user's password
func ChangeUserPassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return ErrInvalidPassword
	}
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return ErrPasswordHash
	}
	user.Password = hashedPassword
	if err := config.DB.Save(&user).Error; err != nil {
		return ErrUserUpdateFailed
	}
	return nil
}
