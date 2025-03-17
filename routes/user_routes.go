package routes

import (
	"azyqs-auth-systems/controllers"
	"azyqs-auth-systems/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterUserRoutes defines routes for user operations
func RegisterUserRoutes(router *mux.Router) {
	protected := router.PathPrefix("/user").Subrouter()
	protected.Use(middlewares.JwtAuthentication)

	protected.HandleFunc("/profile", controllers.ViewProfile).Methods("GET")
	protected.HandleFunc("/profile", controllers.EditProfile).Methods("PUT")
	protected.HandleFunc("/profile", controllers.DeleteProfile).Methods("DELETE")
	protected.HandleFunc("/change-password", controllers.ChangePassword).Methods("PUT")
	protected.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
}
