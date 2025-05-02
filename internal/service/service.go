package service

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/yourusername/rss-reader/internal/config"
	"github.com/yourusername/rss-reader/internal/feed"
	"github.com/yourusername/rss-reader/internal/notification"
	"github.com/yourusername/rss-reader/internal/storage"
)

type Service struct {
	config        *config.Config
	db            *storage.DB
	logger        *log.Logger
	feedManager   *feed.Manager
	notifier      *notification.Notifier
	cronScheduler *cron.Cron
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

func New(cfg *config.Config, db *storage.DB, logger *log.Logger) *Service {
	notifier := notification.New(logger)
	feedManager := feed.NewManager(db, notifier, logger)

	return &Service{
		config:        cfg,
		db:            db,
		logger:        logger,
		feedManager:   feedManager,
		notifier:      notifier,
		cronScheduler: cron.New(),
		stopChan:      make(chan struct{}),
	}
}

func (s *Service) Start() error {
	// Initialize the feed manager
	if err := s.feedManager.Initialize(s.config.Feeds); err != nil {
		return err
	}

	// Start periodic feed updates
	s.cronScheduler.AddFunc("@every "+s.config.RefreshInterval.String(), func() {
		if err := s.feedManager.UpdateAllFeeds(); err != nil {
			s.logger.Printf("Failed to update feeds: %v", err)
		}
	})
	s.cronScheduler.Start()

	// Execute initial feed update
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.feedManager.UpdateAllFeeds(); err != nil {
			s.logger.Printf("Initial feed update failed: %v", err)
		}
	}()

	s.logger.Println("RSS service started successfully")
	return nil
}

func (s *Service) Stop() error {
	s.logger.Println("Stopping RSS service...")

	// Stop the scheduler
	s.cronScheduler.Stop()

	// Signal all goroutines to stop
	close(s.stopChan)

	// Wait for all goroutines to finish
	s.wg.Wait()

	s.logger.Println("RSS service stopped successfully")
	return nil
}
