package main

import (
	"log"

	"tasius.my.id/todolistapi/internal/config"
	"tasius.my.id/todolistapi/internal/infrastructure/db"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	dbConn, err := db.ConnectWithoutMigration(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.Migrate(dbConn); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Migrations completed successfully")
}
