package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// Load file .env, jika ada. Jika tidak, akan lanjut dengan env yang sudah di-set.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file tidak ditemukan, melanjutkan dengan environment variable")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL belum di-set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	log.Println("Berhasil terhubung ke database")
	return DB
}
