package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func StartDB() (*pgxpool.Pool, error) {
	dsn := "postgres://server:production@localhost:5433/forum_db?sslmode=disable"
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = 30
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	// TestDb
	var now time.Time
	err = pool.QueryRow(context.Background(), "SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	fmt.Printf("✅ Connected to Postgres! Server time: %s\n", now)

	return pool, nil
}

// Add a trigger in Postgres that NOTIFY chat, '{"conversation_id":123,"message_id":456}' on new messages. This is essential for live
// func listenChat(ctx context.Context, pool *pgxpool.Pool, handle func(payload string)) error {
// 	conn, err := pool.Acquire(ctx)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	go func() {
// 		defer conn.Release()
// 		if _, err := conn.Exec(ctx, "LISTEN chat"); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		for {
// 			n, err := conn.Conn().WaitForNotification(ctx)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			handle(n.Payload)
// 		}
// 	}()
// 	return nil
// }
