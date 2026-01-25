package main

import (
	"log"
	"os"
	postgresql "social-network/shared/go/postgre"
)

func main() {
	log.Println("Running database migrations...")

	if err := postgresql.RunMigrations(os.Getenv("DATABASE_URL"), os.Getenv("MIGRATE_PATH")); err != nil {
		log.Fatal("Migration failed", err)
	}

	log.Println("Migrations completed successfully.")
	os.Exit(0)
}
