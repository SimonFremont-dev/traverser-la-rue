package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SimonFremont-dev/cv-builder/backend/internal/api"
	"github.com/SimonFremont-dev/cv-builder/backend/internal/config"
	"github.com/SimonFremont-dev/cv-builder/backend/internal/repository"
	"github.com/SimonFremont-dev/cv-builder/backend/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadFromEnv()
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	// Initialize SQLite repository
	sqliteRepo, err := repository.NewSQLiteRepository(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite repository: %v", err)
	}
	defer sqliteRepo.Close()

	// Initialize SQL repository if configured
	var sqlRepo *repository.SQLRepository
	if cfg.SQLDSN != "" {
		sqlRepo, err = repository.NewSQLRepository(cfg.SQLDSN)
		if err != nil {
			log.Fatalf("Failed to initialize SQL repository: %v", err)
		}
		defer sqlRepo.Close()

		// Run migrations
		if err := sqlRepo.RunMigrations("migrations"); err != nil {
			log.Printf("Warning: failed to run migrations: %v", err)
		} else {
			log.Println("SQL migrations applied successfully")
		}
	}

	// Initialize service
	svc := service.NewService(sqliteRepo, sqlRepo)

	// Initialize handlers
	handlers := api.NewHandlers(svc)

	// Create router
	r := mux.NewRouter()
	handlers.SetupRoutes(r)

	// Configure server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
