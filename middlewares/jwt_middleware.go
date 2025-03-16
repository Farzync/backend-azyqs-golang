package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"azyqs-auth-systems/utils"
)

// Response standar
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// writeJSON membantu menulis response dengan format JSON
func writeJSON(w http.ResponseWriter, statusCode int, status, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
	})
}

// JwtAuthentication memeriksa token JWT di header Authorization
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			writeJSON(w, http.StatusForbidden, "error", "Token tidak ditemukan")
			return
		}
		// Format yang umum: "Bearer {token}"
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			writeJSON(w, http.StatusForbidden, "error", "Format token tidak valid")
			return
		}
		tokenPart := splitted[1]
		userID, err := utils.ValidateJWT(tokenPart)
		if err != nil {
			writeJSON(w, http.StatusForbidden, "error", "Token tidak valid")
			return
		}
		// Masukkan userID ke dalam context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
