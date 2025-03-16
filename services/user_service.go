package services

import (
	"azyqs-auth-systems/config"
	serviceErrors "azyqs-auth-systems/errors"
	"azyqs-auth-systems/models"
	"azyqs-auth-systems/utils"
	"errors"

	"gorm.io/gorm"
)

// GetUserByID fetches user data by ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, serviceErrors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile updates username, name, and email for a user
func UpdateUserProfile(userID uint, newUsername, newName, newEmail string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return serviceErrors.ErrUserNotFound
	}

	// Check for username uniqueness if changed
	if newUsername != user.Username {
		var count int64
		config.DB.Model(&models.User{}).
			Where("username = ? AND id != ?", newUsername, userID).
			Count(&count)
		if count > 0 {
			return serviceErrors.ErrUsernameTaken
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
			return serviceErrors.ErrEmailTaken
		}
		user.Email = newEmail
	}

	user.Name = newName

	if err := config.DB.Save(&user).Error; err != nil {
		return serviceErrors.ErrUserUpdateFailed
	}
	return nil
}

// DeleteUser deletes a user after password confirmation
func DeleteUser(userID uint, password string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return serviceErrors.ErrUserNotFound
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return serviceErrors.ErrPasswordMismatch
	}
	if err := config.DB.Delete(&user).Error; err != nil {
		return serviceErrors.ErrUserDeleteFailed
	}
	return nil
}

// ChangeUserPassword changes a user's password
func ChangeUserPassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return serviceErrors.ErrUserNotFound
	}
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return serviceErrors.ErrInvalidPassword
	}
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return serviceErrors.ErrPasswordHash
	}
	user.Password = hashedPassword
	if err := config.DB.Save(&user).Error; err != nil {
		return serviceErrors.ErrUserUpdateFailed
	}
	return nil
}
