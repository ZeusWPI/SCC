// Package lyrics provides a way to work with both synced and plain lyrics
package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/dto"
)

// Lyrics is the common interface for different lyric types
type Lyrics interface {
	GetSong() dto.Song
	Previous(int) []Lyric
	Current() (Lyric, bool)
	Next() (Lyric, bool)
	Upcoming(int) []Lyric
	Progress() float64
}

// Lyric represents a single lyric line.
type Lyric struct {
	Text     string
	Duration time.Duration
}

// New returns a new object that implements the Lyrics interface
func New(song dto.Song) Lyrics {
	// Basic sync
	if song.LyricsType == "synced" {
		return newLRC(song)
	}

	// Plain
	if song.LyricsType == "plain" {
		return newPlain(song)
	}

	// Instrumental
	if song.LyricsType == "instrumental" {
		return newInstrumental(song)
	}

	// No lyrics found
	return newMissing(song)
}
