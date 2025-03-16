package routes

import (
	"azyqs-auth-systems/controllers"
	"azyqs-auth-systems/middlewares"

	"github.com/gorilla/mux"
)

// RegisterRoutes mendefinisikan semua endpoint API
func RegisterRoutes(router *mux.Router) {
	// Endpoint publik
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// Endpoint yang dilindungi JWT
	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middlewares.JwtAuthentication)
	protected.HandleFunc("/profile", controllers.ViewProfile).Methods("GET")
	protected.HandleFunc("/profile", controllers.EditProfile).Methods("PUT")
	protected.HandleFunc("/profile", controllers.DeleteProfile).Methods("DELETE")
	protected.HandleFunc("/change-password", controllers.ChangePassword).Methods("PUT")
}
