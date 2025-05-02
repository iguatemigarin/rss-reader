package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Feed struct {
	URL   string
	Title string
}

type Item struct {
	ID          int64
	FeedURL     string
	Title       string
	Link        string
	PublishedAt time.Time
}

type DB struct {
	db *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	if err := initSchema(db); err != nil {
		db.Close()
		return nil, err
	}

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func initSchema(db *sql.DB) error {
	// Create feeds table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS feeds (
			url TEXT PRIMARY KEY,
			title TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create items table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			feed_url TEXT NOT NULL,
			title TEXT NOT NULL,
			link TEXT NOT NULL,
			published_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (feed_url) REFERENCES feeds(url) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create index for faster queries
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_items_feed_url ON items(feed_url)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) AddFeed(url, title string) error {
	_, err := d.db.Exec("INSERT OR REPLACE INTO feeds (url, title) VALUES (?, ?)", url, title)
	return err
}

func (d *DB) UpdateFeedMetadata(url, title string) error {
	_, err := d.db.Exec("UPDATE feeds SET title = ? WHERE url = ?", title, url)
	return err
}

func (d *DB) RemoveFeed(url string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete items first due to foreign key constraint
	_, err = tx.Exec("DELETE FROM items WHERE feed_url = ?", url)
	if err != nil {
		return err
	}

	// Delete the feed
	_, err = tx.Exec("DELETE FROM feeds WHERE url = ?", url)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *DB) GetAllFeeds() ([]*Feed, error) {
	rows, err := d.db.Query("SELECT url, title FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*Feed
	for rows.Next() {
		feed := &Feed{}
		if err := rows.Scan(&feed.URL, &feed.Title); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, rows.Err()
}

func (d *DB) AddItem(feedURL, title, link string, publishedAt time.Time) error {
	_, err := d.db.Exec(
		"INSERT INTO items (feed_url, title, link, published_at) VALUES (?, ?, ?, ?)",
		feedURL, title, link, publishedAt,
	)
	return err
}

func (d *DB) GetLatestItem(feedURL string) (*Item, error) {
	row := d.db.QueryRow(`
		SELECT id, feed_url, title, link, published_at 
		FROM items 
		WHERE feed_url = ? 
		ORDER BY published_at DESC 
		LIMIT 1
	`, feedURL)

	item := &Item{}
	err := row.Scan(&item.ID, &item.FeedURL, &item.Title, &item.Link, &item.PublishedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (d *DB) GetItems(feedURL string, limit, offset int) ([]*Item, error) {
	rows, err := d.db.Query(`
		SELECT id, feed_url, title, link, published_at 
		FROM items 
		WHERE feed_url = ? 
		ORDER BY published_at DESC 
		LIMIT ? OFFSET ?
	`, feedURL, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		if err := rows.Scan(&item.ID, &item.FeedURL, &item.Title, &item.Link, &item.PublishedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}