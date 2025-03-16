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

// Response standar
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Load file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Tidak menemukan file .env, pastikan environment variable sudah di-set")
	}

	// Argument flag untuk port dengan default 8080
	port := flag.String("port", "", "Port untuk server (default dari .env atau 8080)")
	flag.Parse()

	// Jika tidak ada argumen port, cek dari environment, lalu fallback ke 8080
	finalPort := *port
	if finalPort == "" {
		finalPort = os.Getenv("PORT")
		if finalPort == "" {
			finalPort = "8080"
		}
	}

	// Inisialisasi koneksi database
	db := config.InitDB()
	db.AutoMigrate(&models.User{})

	// Inisialisasi router
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Custom NotFoundHandler untuk route yang tidak ada
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "route_not_found",
		})
	})

	log.Printf("Server berjalan di port %s...\n", finalPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", finalPort), router))
}
