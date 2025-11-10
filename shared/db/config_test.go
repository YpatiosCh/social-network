package db

import (
	"testing"
)

func TestLoadConfigFromEnv_DatabaseURL(t *testing.T) {
	// Set DATABASE_URL and ensure it's parsed
	old := ""
	// t.Setenv is available in Go 1.17+ and will restore env after test
	t.Setenv("DATABASE_URL", "postgres://alice:pwd123@dbhost:5433/mydb?sslmode=require")

	cfg := LoadConfigFromEnv()

	if cfg.Host != "dbhost" {
		t.Fatalf("expected host dbhost, got %s", cfg.Host)
	}
	if cfg.Port != "5433" {
		t.Fatalf("expected port 5433, got %s", cfg.Port)
	}
	if cfg.User != "alice" {
		t.Fatalf("expected user alice, got %s", cfg.User)
	}
	if cfg.Password != "pwd123" {
		t.Fatalf("expected password pwd123, got %s", cfg.Password)
	}
	if cfg.DBName != "mydb" {
		t.Fatalf("expected dbname mydb, got %s", cfg.DBName)
	}
	if cfg.SSLMode != "require" {
		t.Fatalf("expected sslmode require, got %s", cfg.SSLMode)
	}

	_ = old
}

func TestLoadConfigFromEnv_FallbackToDBVars(t *testing.T) {
	// Clear DATABASE_URL and set DB_* vars
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DB_HOST", "fallback-host")
	t.Setenv("DB_PORT", "5434")
	t.Setenv("DB_USER", "bob")
	t.Setenv("DB_PASSWORD", "s3cr3t")
	t.Setenv("DB_NAME", "fallbackdb")
	t.Setenv("SSL_MODE", "disable")

	cfg := LoadConfigFromEnv()

	if cfg.Host != "fallback-host" {
		t.Fatalf("expected host fallback-host, got %s", cfg.Host)
	}
	if cfg.Port != "5434" {
		t.Fatalf("expected port 5434, got %s", cfg.Port)
	}
	if cfg.User != "bob" {
		t.Fatalf("expected user bob, got %s", cfg.User)
	}
	if cfg.Password != "s3cr3t" {
		t.Fatalf("expected password s3cr3t, got %s", cfg.Password)
	}
	if cfg.DBName != "fallbackdb" {
		t.Fatalf("expected dbname fallbackdb, got %s", cfg.DBName)
	}
	if cfg.SSLMode != "disable" {
		t.Fatalf("expected sslmode disable, got %s", cfg.SSLMode)
	}
}
