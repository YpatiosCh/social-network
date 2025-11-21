package db

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfigFromEnv() Config {
	// Support DATABASE_URL (common in docker-compose and 12-factor apps).
	// If present, parse it and populate fields accordingly. Otherwise fall
	// back to individual DB_* environment variables with sensible defaults.
	if raw := os.Getenv("DATABASE_URL"); raw != "" {
		// pg URL may contain query params (eg sslmode)
		u, err := url.Parse(raw)
		if err == nil && u.Scheme != "" {
			host := u.Hostname()
			port := u.Port()
			if port == "" {
				port = "5432"
			}
			user := ""
			pass := ""
			if u.User != nil {
				user = u.User.Username()
				if p, ok := u.User.Password(); ok {
					pass = p
				}
			}
			dbname := strings.TrimPrefix(u.Path, "/")
			// extract sslmode from query if present
			q := u.Query()
			ssl := q.Get("sslmode")
			if ssl == "" {
				ssl = getEnv("SSL_MODE", "disable")
			}

			return Config{
				Host:     host,
				Port:     port,
				User:     user,
				Password: pass,
				DBName:   dbname,
				SSLMode:  ssl,
			}
		}
		// if parsing failed, fall through to env vars below
	}

	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "secret"),
		DBName:   getEnv("DB_NAME", "social_users"),
		SSLMode:  getEnv("SSL_MODE", "disable"),
	}
}

func (c Config) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
