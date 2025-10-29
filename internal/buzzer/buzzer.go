// Package buzzer provides all interactions with the buzzer
package buzzer

import (
	"os/exec"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Client represents a buzzer
type Client struct {
	Song []string
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
	return &Client{
		Song: config.GetDefaultStringSlice("backend.buzzer.song", defaultSong),
	}
}

// Play plays the buzzer
func (c *Client) Play() {
	// See `man beep` for more information
	cmd := exec.Command("beep", c.Song...)
	err := cmd.Run()
	if err != nil {
		zap.S().Error("Error running command 'beep' %v", err)
	}
}
