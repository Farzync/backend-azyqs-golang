package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection using environment variables.
func InitDB() *gorm.DB {
	// Load .env file if it exists, otherwise continue with existing environment variables.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, continuing with environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Error: DATABASE_URL is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error: failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully")
	return DB
}
