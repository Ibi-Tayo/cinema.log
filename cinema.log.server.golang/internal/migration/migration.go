package migration

import (
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed goose/*.sql
var embedMigrations embed.FS

// RunMigrations runs all pending migrations using Goose
func RunMigrations(db *sql.DB) error {
	log.Println("Starting database migrations...")
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// check current version
	currentVersion, err := goose.GetDBVersion(db)
	if err != nil {
		log.Printf("Warning: Could not get current DB version (expected on first run): %v\n", err)
		currentVersion = 0
	} else {
		log.Printf("Current migration version: %d\n", currentVersion)
	}

	// run migrations
	log.Println("Running migrations...")
	if err := goose.Up(db, "goose"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// check if new version is different from current version
	newVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get new version: %w", err)
	}

	if newVersion == currentVersion {
		log.Println("No new migrations applied, database up to date")
	} else {
		log.Printf("Migrations applied successfully: %d -> %d\n", currentVersion, newVersion)
	}

	return nil
}

// RollbackMigration rolls back one migration
func RollbackMigration(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Down(db, "goose"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Status(db, "goose"); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

// ResetDatabase resets the database by rolling back all migrations
func ResetDatabase(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Reset(db, "goose"); err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	return nil
}
