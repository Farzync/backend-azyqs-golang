package routes

import (
	"azyqs-auth-systems/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterAuthRoutes defines routes for authentication
func RegisterAuthRoutes(router *mux.Router) {
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", controllers.Register).Methods("POST")
	authRouter.HandleFunc("/login", controllers.Login).Methods("POST")
	authRouter.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
}
