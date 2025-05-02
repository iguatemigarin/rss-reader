package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/rss-reader/internal/config"
	"github.com/yourusername/rss-reader/internal/service"
	"github.com/yourusername/rss-reader/internal/storage"
)

func main() {
	logger := log.New(os.Stdout, "RSS-READER: ", log.LstdFlags)
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := storage.NewDB(cfg.DatabasePath)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create and start the service
	svc := service.New(cfg, db, logger)
	if err := svc.Start(); err != nil {
		logger.Fatalf("Failed to start service: %v", err)
	}

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Clean shutdown
	logger.Println("Shutting down service...")
	if err := svc.Stop(); err != nil {
		logger.Printf("Error during shutdown: %v", err)
	}
}