package database

import (
	"log"

	"github.com/CodeEnthusiast09/country-currency-api/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

// New connects to the database and returns a Database instance.
//
// NOTE ON MIGRATIONS:
// Migrations are intentionally NOT run here in production.
// They are handled by entrypoint.sh (via `atlas migrate apply`) BEFORE
// this application process starts. This is the industry standard pattern:
//
//	entrypoint.sh:
//	  1. Wait for DB to be ready
//	  2. atlas migrate apply   ← migrations happen here
//	  3. exec /app/server      ← only then does this code run
//
// This separation means:
//   - If a migration fails, the app never starts (clear failure signal)
//   - Multiple app replicas won't race to migrate the same DB simultaneously
//   - Migration concerns are kept out of application code
//
// For local development, run migrations manually:
//
//	atlas migrate apply --env gorm --url "mysql://user:pass@localhost:3306/db"
func New(cfg *config.Config) *Database {
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	return &Database{DB: db}
}
