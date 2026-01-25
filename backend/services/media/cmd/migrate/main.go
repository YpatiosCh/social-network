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

	if err := postgresql.RunMigrations(os.Getenv("DATABASE_URL"), os.Getenv("MIGRATE_PATH")); err != nil {
		tele.Fatal("migration failed, erro: " + err.Error())
	}

	tele.Info(ctx, "Migrations completed successfully.")
	os.Exit(0)
}
