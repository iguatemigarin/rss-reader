package ui

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/yourusername/rss-reader/internal/config"
)

func OnReady(cfg *config.Config, logger *log.Logger) {
	systray.SetIcon(getIcon())
	systray.SetTitle("RSS Reader")
	systray.SetTooltip("RSS Reader")

	mFeeds := systray.AddMenuItem("Manage Feeds", "Open the feed management UI")
	mSettings := systray.AddMenuItem("Settings", "Configure RSS Reader")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit RSS Reader")

	go func() {
		for {
			select {
			case <-mFeeds.ClickedCh:
				logger.Println("Opening feed management UI")
				openURL("http://localhost:3000/feeds")
			case <-mSettings.ClickedCh:
				logger.Println("Opening settings")
				openURL("http://localhost:3000/settings")
			case <-mQuit.ClickedCh:
				logger.Println("User requested quit")
				systray.Quit()
				return
			}
		}
	}()
}

func OnExit(logger *log.Logger) {
	logger.Println("Exiting UI")
}

func openURL(url string) {
	var cmd *exec.Cmd

switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}

cmd.Start()
}

func getIcon() []byte {
	// This is a placeholder for an actual icon
	// In a real implementation, you would load an .ico file
	return []byte{
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x10, 0x10,
		0x00, 0x00, 0x01, 0x00, 0x08, 0x00, 0x68, 0x05,
		0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x28, 0x00,
		0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x20, 0x00,
	}
}

