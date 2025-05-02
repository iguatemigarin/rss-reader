package notification

import (
	"log"
	"os/exec"
)

type Notifier struct {
	logger *log.Logger
}

func New(logger *log.Logger) *Notifier {
	return &Notifier{
		logger: logger,
	}
}

func (n *Notifier) NotifyNewItem(feedTitle, itemTitle, itemLink string) {
	n.logger.Printf("New item: %s - %s", feedTitle, itemTitle)
	
	// Use AppleScript to display a native macOS notification
	script := `display notification "` + itemTitle + `" with title "` + feedTitle + `"`
	
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		n.logger.Printf("Failed to send notification: %v", err)
	}
}