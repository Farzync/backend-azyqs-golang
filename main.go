package main

import (
	"encoding/json"
	"log"
	"net/http"

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
			Message: "route tidak ditemukan",
		})
	})

	log.Println("Server berjalan di port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
