// Package buzzer provides all interactions with the buzzer
package buzzer

import (
	"os/exec"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Client represents a buzzer
type Client struct {
	hasBuzzer bool
	song      []string
}

var defaultSong = []string{
	"-n", "-f880", "-l100", "-d0",
	"-n", "-f988", "-l100", "-d0",
	"-n", "-f588", "-l100", "-d0",
	"-n", "-f989", "-l100", "-d0",
	"-n", "-f660", "-l200", "-d0",
	"-n", "-f660", "-l200", "-d0",
	"-n", "-f588", "-l100", "-d0",
	"-n", "-f555", "-l100", "-d0",
	"-n", "-f495", "-l100", "-d0",
}

// New returns a new buzzer instance
func New() *Client {
	hasBuzzer := false
	if _, err := exec.LookPath("beep"); err == nil {
		hasBuzzer = true
	}

	if !hasBuzzer {
		zap.S().Debug("No beep executable found.\nMock messages will be used instead")
	}

	return &Client{
		hasBuzzer: hasBuzzer,
		song:      config.GetDefaultStringSlice("backend.buzzer.song", defaultSong),
	}
}

// Play plays the buzzer
func (c *Client) Play() {
	if !c.hasBuzzer {
		zap.S().Info("BEEEEEEEP")
		return
	}

	// See `man beep` for more information
	cmd := exec.Command("beep", c.song...)
	err := cmd.Run()
	if err != nil {
		zap.S().Error("Error running command 'beep' %v", err)
	}
}
