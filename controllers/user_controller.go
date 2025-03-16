package controllers

import (
	"azyqs-auth-systems/middlewares"
	"azyqs-auth-systems/services"
	"encoding/json"
	"net/http"
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

// Register: POST /auth/register
func Register(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "invalid_input", nil)
		return
	}

	err := services.RegisterUser(userInput.Username, userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, "success", "registration_successful", nil)
}

// Login: POST /auth/login
func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "invalid_input", nil)
		return
	}

	token, err := services.LoginUser(input.Username, input.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, "error", err.Error(), nil)
		return
	}

	writeJSON(w, http.StatusOK, "success", "login_successful", map[string]string{"token": token})
}

// View Profile: GET /user/profile
func ViewProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(uint)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", "user_id_not_found", nil)
		return
	}

	user, err := services.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, "error", "user_not_found", nil)
		return
	}

	writeJSON(w, http.StatusOK, "success", "profile_found", user)
}

// Edit Profile: PUT /user/profile
func EditProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(uint)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", "user_id_not_found", nil)
		return
	}

	var input struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "invalid_input", nil)
		return
	}

	err := services.UpdateUserProfile(userID, input.Username, input.Name, input.Email)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, "success", "profile_updated", nil)
}

// Delete Profile: DELETE /user/profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(uint)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", "user_id_not_found", nil)
		return
	}

	var input struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "invalid_input", nil)
		return
	}

	err := services.DeleteUser(userID, input.Password)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, "success", "profile_deleted", nil)
}

// Change Password: PUT /user/change-password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey).(uint)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, "error", "user_id_not_found", nil)
		return
	}

	var input struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", "invalid_input", nil)
		return
	}

	if input.NewPassword != input.ConfirmNewPassword {
		writeJSON(w, http.StatusBadRequest, "error", "password_mismatch", nil)
		return
	}

	err := services.ChangeUserPassword(userID, input.OldPassword, input.NewPassword)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	writeJSON(w, http.StatusOK, "success", "password_changed", nil)
}
