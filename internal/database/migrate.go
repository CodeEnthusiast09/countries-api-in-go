package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending "up" migrations from the migrations/ folder.
// It receives a *sql.DB (standard library) — NOT a *gorm.DB.
// golang-migrate works at a lower level than GORM, so it needs the raw driver.
func RunMigrations(sqlDB *sql.DB, dbName string) error {
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{
		DatabaseName: dbName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	// "file://migrations" tells migrate to look for .sql files in the migrations/ folder
	// relative to where you run the binary (i.e., your project root)
	m, err := migrate.NewWithDatabaseInstance("file://migrations", dbName, driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// ErrNoChange just means "nothing new to run" — that's fine, not a real error
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations ran successfully")
	return nil
}
