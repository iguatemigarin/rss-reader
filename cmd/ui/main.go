package main

import (
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/yourusername/rss-reader/internal/config"
	"github.com/yourusername/rss-reader/internal/ui"
)

func main() {
	logger := log.New(os.Stdout, "RSS-UI: ", log.LstdFlags)
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Start the system tray application
	systray.Run(
		func() { ui.OnReady(cfg, logger) },
		func() { ui.OnExit(logger) },
	)
}