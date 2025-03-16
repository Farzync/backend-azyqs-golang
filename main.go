package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"azyqs-auth-systems/config"
	"azyqs-auth-systems/models"
	"azyqs-auth-systems/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Standard response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, make sure environment variables are set")
	}

	// Command-line flag for server port (default 8080)
	port := flag.String("port", "", "Server port (default from .env or 8080)")
	flag.Parse()

	// If no port argument is provided, check environment variables, fallback to 8080
	finalPort := *port
	if finalPort == "" {
		finalPort = os.Getenv("PORT")
		if finalPort == "" {
			finalPort = "8080"
		}
	}

	// Initialize database connection
	db := config.InitDB()
	db.AutoMigrate(&models.User{})

	// Initialize router
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Custom NotFoundHandler for non-existent routes
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "route_not_found",
		})
	})

	log.Printf("Server is running on port %s...\n", finalPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", finalPort), router))
}
