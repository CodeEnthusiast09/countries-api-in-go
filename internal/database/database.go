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

// FOR AUTO-MiGRATION
// func New(cfg *config.Config) *Database {
// 	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
//
// 	if err := db.AutoMigrate(&models.Country{}); err != nil {
// 		log.Fatalf("AutoMigrate failed: %v", err)
// 	}
//
// 	log.Println("Database connected and migrated")
// 	return &Database{DB: db}
// }

func New(cfg *config.Config) *Database {
	db, err := gorm.Open(mysql.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 👇👇👇 - for go-migrate
	// gorm.DB wraps the standard library *sql.DB.
	// We unwrap it here because golang-migrate needs the raw driver.
	// sqlDB, err := db.DB()
	// if err != nil {
	// 	log.Fatalf("Failed to get underlying sql.DB: %v", err)
	// }
	//
	// if err := RunMigrations(sqlDB, cfg.DBName); err != nil {
	// 	log.Fatalf("Migration error: %v", err)
	// }
	// 👆👆👆 - for go-migrate

	// 👇👇👇 - for Atlas
	// Pass the folder path (not a file:// URL) — os.DirFS in RunMigrations handles it
	if err := RunMigrations("migrations", cfg.AtlasDatabaseURL); err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	// 👆👆👆 - for Atlas

	log.Println("Database connected and migrated")
	return &Database{DB: db}
}
