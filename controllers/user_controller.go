package controllers

import (
	"azyqs-auth-systems/services"
	"encoding/json"
	"net/http"
)

// Response standar
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Register: POST /register
func Register(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Input tidak valid"})
		return
	}

	err := services.RegisterUser(userInput.Username, userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Registrasi berhasil"})
}

// Login: POST /login
func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Input tidak valid"})
		return
	}

	token, err := services.LoginUser(input.Username, input.Password)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Login berhasil",
		Data:    map[string]string{"token": token},
	})
}

// View Profile: GET /profile
func ViewProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Profile ditemukan", Data: user})
}

// Edit Profile: PUT /profile
func EditProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	var input struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Input tidak valid"})
		return
	}
	err := services.UpdateUserProfile(userID.(uint), input.Username, input.Name, input.Email)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Profile berhasil diupdate"})
}


// Delete Profile: DELETE /profile
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	var input struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Input tidak valid"})
		return
	}
	err := services.DeleteUser(userID.(uint), input.Password)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Profile berhasil dihapus"})
}

// Change Password: PUT /change-password
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	var input struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Input tidak valid"})
		return
	}
	if input.NewPassword != input.ConfirmNewPassword {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Password baru tidak cocok"})
		return
	}
	err := services.ChangeUserPassword(userID.(uint), input.OldPassword, input.NewPassword)
	if err != nil {
		json.NewEncoder(w).Encode(Response{Status: "error", Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Password berhasil diubah"})
}
