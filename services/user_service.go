package services

import (
	"azyqs-auth-systems/config"
	"azyqs-auth-systems/models"
	"azyqs-auth-systems/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

// RegisterUser: Registrasi user baru
func RegisterUser(username, name, email, password string) error {
	// Cek apakah username/email sudah terpakai
	var count int64
	config.DB.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return errors.New("username atau email sudah digunakan")
	}
	
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	
	user := models.User{
		Username: username,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}
	
	if err := config.DB.Create(&user).Error; err != nil {
		// Tangani duplicate error dari database
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return errors.New("username atau email sudah digunakan")
		}
		return err
	}
	return nil
}

// LoginUser: Autentikasi user dan mengembalikan JWT token
func LoginUser(username, password string) (string, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("user tidak ditemukan")
		}
		return "", err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("password salah")
	}
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetUserByID: Mengambil data user berdasarkan ID
func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile: Update username, name, dan email user
func UpdateUserProfile(userID uint, newUsername, newName, newEmail string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}

	// Jika username diubah, periksa apakah sudah ada yang pakai (selain user ini)
	if newUsername != user.Username {
		var count int64
		config.DB.Model(&models.User{}).
			Where("username = ? AND id != ?", newUsername, userID).
			Count(&count)
		if count > 0 {
			return errors.New("username sudah digunakan")
		}
		user.Username = newUsername
	}

	// Jika email diubah, periksa apakah sudah ada yang pakai (selain user ini)
	if newEmail != user.Email {
		var count int64
		config.DB.Model(&models.User{}).
			Where("email = ? AND id != ?", newEmail, userID).
			Count(&count)
		if count > 0 {
			return errors.New("email sudah digunakan")
		}
		user.Email = newEmail
	}

	// Update nama
	user.Name = newName

	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUser: Hapus user dengan konfirmasi password
func DeleteUser(userID uint, password string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return errors.New("password tidak cocok")
	}
	if err := config.DB.Delete(&user).Error; err != nil {
		return err
	}
	return nil
}

// ChangeUserPassword: Ubah password user
func ChangeUserPassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("password lama salah")
	}
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
