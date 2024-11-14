// Package buzzer provides all interactions with the buzzer
package buzzer

import (
	"os/exec"

	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

// Buzzer represents a buzzer
type Buzzer struct {
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
func New() *Buzzer {
	song := config.GetDefaultStringSlice("buzzer.song", defaultSong)
	return &Buzzer{
		Song: song,
	}
}

// Play plays the buzzer
func (b *Buzzer) Play() {
	// See `man beep` for more information
	cmd := exec.Command("beep", b.Song...)
	err := cmd.Run()

	if err != nil {
		zap.L().Error("Error running command 'beep'", zap.Error(err))
	}
}
