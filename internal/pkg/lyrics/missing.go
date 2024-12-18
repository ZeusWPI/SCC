package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/dto"
)

// Missing represents lyrics that are absent
type Missing struct {
	song   dto.Song
	lyrics Lyric
	given  bool
}

func newMissing(song dto.Song) Lyrics {
	lyric := Lyric{
		Text:     "Missing lyrics\n\nHelp the open source community by adding them to\nhttps://lrclib.net/",
		Duration: time.Duration(song.DurationMS) * time.Millisecond,
	}

	return &Missing{song: song, lyrics: lyric, given: false}
}

// GetSong returns the song associated to the lyrics
func (m *Missing) GetSong() dto.Song {
	return m.song
}

// Previous provides the previous `amount` of lyrics without affecting the current lyric
// In this case it's alway nothing
func (m *Missing) Previous(_ int) []Lyric {
	return []Lyric{}
}

// Current provides the current lyric if any.
// If the song is finished the boolean is set to false
func (m *Missing) Current() (Lyric, bool) {
	if m.given {
		return Lyric{}, false
	}

	return m.lyrics, true
}

// Next provides the next lyric.
// If the lyrics are finished the boolean is set to false
func (m *Missing) Next() (Lyric, bool) {
	if m.given {
		return Lyric{}, false
	}

	m.given = true

	return m.lyrics, true
}

// Upcoming provides the next `amount` lyrics without affecting the current lyric
// In this case it's always empty
func (m *Missing) Upcoming(_ int) []Lyric {
	return []Lyric{}
}

// Progress shows the fraction of lyrics that have been used.
func (m *Missing) Progress() float64 {
	if m.given {
		return 1
	}

	return 0
}
