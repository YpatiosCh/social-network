package main

import (
	"log"
	"os"

	"social-network/shared/db"
)

func main() {
	cfg := db.LoadConfigFromEnv()

	log.Println("Running database migrations...")

	if err := db.RunMigrations(cfg, "./migrations"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
