package database

import (
	"log"

	"github.com/CodeEnthusiast09/country-currency-api/internal/config"
	"github.com/CodeEnthusiast09/country-currency-api/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New(cfg *config.Config) *Database {
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&models.Country{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	log.Println("Database connected and migrated")
	return &Database{DB: db}
}
