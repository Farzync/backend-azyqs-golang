package routes

import (
	"azyqs-auth-systems/controllers"
	"azyqs-auth-systems/middlewares"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// ErrorResponse defines the standard error structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// methodNotAllowedHandler handles disallowed HTTP methods
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Message: "method_not_allowed",
	})
}

// RegisterRoutes defines all API endpoints
func RegisterRoutes(router *mux.Router) {
	// Public auth endpoints
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", controllers.Register).Methods("POST")
	authRouter.HandleFunc("/login", controllers.Login).Methods("POST")
	authRouter.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	// Protected user profile endpoints
	protected := router.PathPrefix("/user").Subrouter()
	protected.Use(middlewares.JwtAuthentication)
	protected.HandleFunc("/profile", controllers.ViewProfile).Methods("GET")
	protected.HandleFunc("/profile", controllers.EditProfile).Methods("PUT")
	protected.HandleFunc("/profile", controllers.DeleteProfile).Methods("DELETE")
	protected.HandleFunc("/change-password", controllers.ChangePassword).Methods("PUT")
	protected.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
}
