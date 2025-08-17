package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"cinema.log.server.golang/internal/database"
	"cinema.log.server.golang/internal/migration"
)

func main() {
	var command = flag.String("command", "", "Migration command: up, down, status, reset")
	flag.Parse()

	if *command == "" {
		fmt.Println("Usage: go run cmd/migrate/main.go -command=<up|down|status|reset>")
		os.Exit(1)
	}

	// Initialize database connection
	db := database.New()
	defer db.Close()

	switch *command {
	case "up":
		if err := migration.RunMigrations(db); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if err := migration.RollbackMigration(db); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		fmt.Println("Migration rolled back successfully")

	case "status":
		if err := migration.GetMigrationStatus(db); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

	case "reset":
		if err := migration.ResetDatabase(db); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		fmt.Println("Database reset successfully")

	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: up, down, status, reset")
		os.Exit(1)
	}
}
