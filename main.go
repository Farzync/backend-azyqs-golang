package main

import (
	"flag"
	"fmt"
	"log"
	"net"
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

func isPortAvailable(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

func main() {
	log.Printf("Starting server...")

	// Load .env file if available
	log.Printf("Loading environment variables...")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, make sure environment variables are set")
	}

	// Command-line flag for server port (default 8080)
	log.Printf("Parsing command-line arguments...")
	port := flag.String("port", "", "Server port (default from .env or 8080)")
	flag.Parse()

	// Determine final port
	log.Printf("Checking server port...")
	finalPort := *port
	if finalPort == "" {
		finalPort = os.Getenv("PORT")
		if finalPort == "" {
			finalPort = "8080"
		}
	}

	// Check if the port is already in use
	if !isPortAvailable(finalPort) {
		log.Fatalf("Port %s is already in use. Please choose a different port.", finalPort)
	}

	// Initialize database connection
	log.Printf("Initializing database connection...")
	db := config.InitDB()
	db.AutoMigrate(&models.User{})

	// Initialize router
	log.Printf("Initializing router...")
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	log.Printf("Server is up and running on port %s...\n", finalPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", finalPort), router))
}
