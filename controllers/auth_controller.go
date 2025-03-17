package controllers

import (
	"azyqs-auth-systems/errors"
	"azyqs-auth-systems/services"
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidInput.Error(), nil)
		return
	}

	err := services.RegisterUser(userInput.Username, userInput.Name, userInput.Email, userInput.Password)
	if err != nil {
		switch err {
		case errors.ErrDuplicateRecord:
			writeJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
		case errors.ErrPasswordHash:
			writeJSON(w, http.StatusInternalServerError, "error", errors.ErrInternalServer.Error(), nil)
		default:
			writeJSON(w, http.StatusInternalServerError, "error", errors.ErrInternalServer.Error(), nil)
		}
		return
	}

	writeJSON(w, http.StatusOK, "success", "registration_successful", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, "error", errors.ErrInvalidInput.Error(), nil)
		return
	}

	token, err := services.LoginUser(input.Username, input.Password)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound, errors.ErrInvalidPassword:
			writeJSON(w, http.StatusUnauthorized, "error", errors.ErrUnauthorized.Error(), nil)
		default:
			writeJSON(w, http.StatusInternalServerError, "error", errors.ErrInternalServer.Error(), nil)
		}
		return
	}

	writeJSON(w, http.StatusOK, "success", "login_successful", map[string]string{"token": token})
}
