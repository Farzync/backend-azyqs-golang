package controllers

import (
	"azyqs-auth-systems/errors"
	"azyqs-auth-systems/middlewares"
	"azyqs-auth-systems/services"
	"azyqs-auth-systems/validators"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, status, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// View Profile: GET /user/profile
func ViewProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", errors.ErrUserIDNotFound.Error(), nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidUserID.Error(), nil)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, "error", errors.ErrUserNotFound.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, "success", "profile_found", user)
}

// Edit Profile: PUT /user/profile
func EditProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", errors.ErrUserIDNotFound.Error(), nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidUserID.Error(), nil)
		return
	}

	var input struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidInput.Error(), nil)
		return
	}

	// Validasi input pengguna
	if input.Username != "" {
		if err := validators.ValidateUsername(input.Username); err != nil {
			writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
			return
		}
	}

	if input.Name != "" {
		if err := validators.ValidateName(input.Name); err != nil {
			writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
			return
		}
	}

	if input.Email != "" {
		if err := validators.ValidateEmail(input.Email); err != nil {
			writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
			return
		}
	}

	err = services.UpdateUserProfile(userID, input.Username, input.Name, input.Email)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, "success", "profile_updated", nil)
}

// Delete Profile: DELETE /user/profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", errors.ErrUserIDNotFound.Error(), nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidUserID.Error(), nil)
		return
	}

	var input struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidInput.Error(), nil)
		return
	}

	// Validasi password
	if err := validators.ValidatePassword(input.Password); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	err = services.DeleteUser(userID, input.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, "success", "profile_deleted", nil)
}

// Change Password: PUT /user/change-password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", errors.ErrUserIDNotFound.Error(), nil)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidUserID.Error(), nil)
		return
	}

	var input struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidInput.Error(), nil)
		return
	}

	// Validasi password
	if err := validators.ValidatePassword(input.OldPassword); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	if err := validators.ValidatePassword(input.NewPassword); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	if input.NewPassword != input.ConfirmNewPassword {
		writeJSON(w, http.StatusBadRequest, "error", "password_mismatch", nil)
		return
	}

	err = services.ChangeUserPassword(userID, input.OldPassword, input.NewPassword)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, "success", "password_changed", nil)
}