package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"os"
	"time"
)

type Config struct {
	// Database configuration
	DatabasePath string
	
	// Feed configuration
	RefreshInterval time.Duration
	Feeds           []string
	
	// Notification configuration
	NotificationsEnabled bool
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	
	configDir := filepath.Join(homeDir, ".rss-reader")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	
	// Set defaults
	viper.SetDefault("database.path", filepath.Join(configDir, "feeds.db"))
	viper.SetDefault("feeds.refresh_interval", 5*time.Minute)
	viper.SetDefault("feeds.sources", []string{})
	viper.SetDefault("notifications.enabled", true)
	
	// Create default config if it doesn't exist
	configPath := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			return nil, err
		}
	}
	
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	
	return &Config{
		DatabasePath:         viper.GetString("database.path"),
		RefreshInterval:      viper.GetDuration("feeds.refresh_interval"),
		Feeds:                viper.GetStringSlice("feeds.sources"),
		NotificationsEnabled: viper.GetBool("notifications.enabled"),
	}, nil
}

func (c *Config) Save() error {
	viper.Set("database.path", c.DatabasePath)
	viper.Set("feeds.refresh_interval", c.RefreshInterval)
	viper.Set("feeds.sources", c.Feeds)
	viper.Set("notifications.enabled", c.NotificationsEnabled)
	
	return viper.WriteConfig()
}