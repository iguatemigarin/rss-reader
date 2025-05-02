package feed

import (
	"log"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/yourusername/rss-reader/internal/notification"
	"github.com/yourusername/rss-reader/internal/storage"
)

type Manager struct {
	db         *storage.DB
	notifier   *notification.Notifier
	logger     *log.Logger
	parser     *gofeed.Parser
	feedUrls   []string
	mutex      sync.RWMutex
}

func NewManager(db *storage.DB, notifier *notification.Notifier, logger *log.Logger) *Manager {
	return &Manager{
		db:       db,
		notifier: notifier,
		logger:   logger,
		parser:   gofeed.NewParser(),
		feedUrls: make([]string, 0),
		mutex:    sync.RWMutex{},
	}
}

func (m *Manager) Initialize(feedUrls []string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.feedUrls = feedUrls
	
	// Load feeds from database
	storedFeeds, err := m.db.GetAllFeeds()
	if err != nil {
		return err
	}
	
	// Add any feeds that are in config but not in the database
	for _, url := range feedUrls {
		found := false
		for _, storedFeed := range storedFeeds {
			if url == storedFeed.URL {
				found = true
				break
			}
		}
		
		if !found {
			if err := m.db.AddFeed(url, ""); err != nil {
				m.logger.Printf("Failed to add feed %s: %v", url, err)
			}
		}
	}
	
	return nil
}

func (m *Manager) UpdateAllFeeds() error {
	m.mutex.RLock()
	urls := make([]string, len(m.feedUrls))
	copy(urls, m.feedUrls)
	m.mutex.RUnlock()
	
	var wg sync.WaitGroup
	errorCh := make(chan error, len(urls))
	
	for _, url := range urls {
		wg.Add(1)
		go func(feedURL string) {
			defer wg.Done()
			
			if err := m.updateFeed(feedURL); err != nil {
				errorCh <- err
			}
		}(url)
	}
	
	wg.Wait()
	close(errorCh)
	
	// Collect errors
	var errs []error
	for err := range errorCh {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		m.logger.Printf("Encountered %d errors while updating feeds", len(errs))
		return errs[0]
	}
	
	return nil
}

func (m *Manager) updateFeed(url string) error {
	m.logger.Printf("Updating feed: %s", url)
	
	feed, err := m.parser.ParseURL(url)
	if err != nil {
		m.logger.Printf("Failed to parse feed %s: %v", url, err)
		return err
	}
	
	// Update feed metadata
	if err := m.db.UpdateFeedMetadata(url, feed.Title); err != nil {
		return err
	}
	
	// Get latest stored item for this feed
	latestItem, err := m.db.GetLatestItem(url)
	if err != nil {
		return err
	}
	
	var latestDate time.Time
	if latestItem != nil {
		latestDate = latestItem.PublishedAt
	}
	
	// Check for new items
	var newItems []*gofeed.Item
	for _, item := range feed.Items {
		pubDate := item.PublishedParsed
		if pubDate == nil {
			continue
		}
		
		if pubDate.After(latestDate) {
			newItems = append(newItems, item)
		}
	}
	
	// Add new items to the database and send notifications
	for _, item := range newItems {
		if err := m.db.AddItem(url, item.Title, item.Link, *item.PublishedParsed); err != nil {
			m.logger.Printf("Failed to add item %s: %v", item.Title, err)
			continue
		}
		
		// Send notification for new item
		m.notifier.NotifyNewItem(feed.Title, item.Title, item.Link)
	}
	
	m.logger.Printf("Updated feed %s, found %d new items", url, len(newItems))
	return nil
}

func (m *Manager) AddFeed(url string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Check if feed already exists
	for _, existingURL := range m.feedUrls {
		if existingURL == url {
			return nil
		}
	}
	
	// Try to parse the feed to validate it
	feed, err := m.parser.ParseURL(url)
	if err != nil {
		return err
	}
	
	// Add to database
	if err := m.db.AddFeed(url, feed.Title); err != nil {
		return err
	}
	
	// Add to in-memory list
	m.feedUrls = append(m.feedUrls, url)
	
	// Update the feed to get initial items
	return m.updateFeed(url)
}

func (m *Manager) RemoveFeed(url string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Remove from database
	if err := m.db.RemoveFeed(url); err != nil {
		return err
	}
	
	// Remove from in-memory list
	for i, existingURL := range m.feedUrls {
		if existingURL == url {
			m.feedUrls = append(m.feedUrls[:i], m.feedUrls[i+1:]...)
			break
		}
	}
	
	return nil
}