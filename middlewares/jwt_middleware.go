package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"azyqs-auth-systems/utils"
)

// ErrorResponse defines the standard error response structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// writeJSON helps write a JSON response with a consistent format
func writeJSON(w http.ResponseWriter, statusCode int, status, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  status,
		Message: message,
	})
}

// Define a custom type for the context key
type contextKey string

const UserIDKey contextKey = "userID"

// JwtAuthentication validates the JWT token in the Authorization header
func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		
		if tokenHeader == "" {
			writeJSON(w, http.StatusForbidden, "error", "token_not_found")
			return
		}

		// Common format: "Bearer {token}"
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			writeJSON(w, http.StatusForbidden, "error", "token_invalid_format")
			return
		}

		tokenPart := splitted[1]
		userID, err := utils.ValidateJWT(tokenPart)
		if err != nil {
			switch err.Error() {
			case "token_expired":
				writeJSON(w, http.StatusForbidden, "error", "token_expired")
			case "token_invalid_signature":
				writeJSON(w, http.StatusForbidden, "error", "token_invalid_signature")
			default:
				writeJSON(w, http.StatusForbidden, "error", "token_invalid")
			}
			return
		}

		// Pass the userID into the request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
