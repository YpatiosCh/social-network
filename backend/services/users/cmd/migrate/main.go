package main

import (
	"context"
	"os"

	postgresql "social-network/shared/go/postgre"
	tele "social-network/shared/go/telemetry"
)

func main() {
	ctx := context.Background()
	tele.Info(ctx, "Running database migrations...")
	dbUrl := os.Getenv("DATABASE_URL")
	if err := postgresql.RunMigrations(dbUrl, os.Getenv("MIGRATE_PATH")); err != nil {
		tele.Fatalf("Migration failed %s", err.Error())
	}

	tele.Info(ctx, "Migrations completed successfully.")
	os.Exit(0)
}
