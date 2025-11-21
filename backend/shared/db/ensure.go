package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func EnsureDatabaseExists(cfg *Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("connect to postgres failed: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check database existence failed: %w", err)
	}

	if !exists {
		_, err := db.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.DBName))
		if err != nil {
			return fmt.Errorf("create database failed: %w", err)
		}
		fmt.Printf("✅ Created database %s\n", cfg.DBName)
	} else {
		fmt.Printf("ℹ️ Database %s already exists\n", cfg.DBName)
	}

	return nil
}
