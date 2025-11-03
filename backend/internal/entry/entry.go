package entry

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/database"
	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/handlers"
)

// server starting sequence
func Start() {
	// set database
	dbPool, err := database.StartDB()
	if err != nil {
		log.Fatal(err)
	}
	db := database.Database{Pool: dbPool}

	// set handlers
	handlers := handlers.Handlers{}
	handlers.Db = &db

	// set server
	var server http.Server
	server.Handler = handlers.SetHandlers()
	server.Addr = "localhost:8081"

	go func() {
		log.Printf("Server running on https://%s\n", server.Addr)
		if err := server.ListenAndServeTLS("localhost+2.pem", "localhost+2-key.pem"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServeTLS failed: %v", err)
		}
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful server Shutdown Failed: %v", err)
	}
	log.Println("Server stopped")
}
