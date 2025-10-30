// Package lyrics provides a way to work with both synced and plain lyrics
package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

// Lyrics is the common interface for different lyric types
type Lyrics interface {
	GetSong() model.Song    // GetSong returns the song associated to the lyrics
	Previous(int) []Lyric   // Previous returns the previous 'int amount' of lyrics without affecting the current lyric
	Current() (Lyric, bool) // Current provides the current lyric and a bool indicating if there is one
	Next() (Lyric, bool)    // Next returns the next lyric and a bool if there are any
	Upcoming(int) []Lyric   // Upcoming returns the next `int amount` of lyrics without affecting the current lyric
	Progress() float64      // Progress returns the fraction of lyrics that have been used
}

// Lyric represents a single lyric line.
type Lyric struct {
	Text     string
	Duration time.Duration
}

// New returns a new object that implements the Lyrics interface
func New(song model.Song) Lyrics {
	switch song.LyricsType {
	case model.LyricsSynced:
		return newLRC(song)
	case model.LyricsInstrumental:
		return newInstrumental(song)
	case model.LyricsPlain:
		return newPlain(song)
	default:
		return newMissing(song)
	}
}
