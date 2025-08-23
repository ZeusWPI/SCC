package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

// Missing represents lyrics that are absent
type Missing struct {
	song   model.Song
	lyrics Lyric
}

func newMissing(song model.Song) Lyrics {
	lyric := Lyric{
		Text:     "Missing lyrics\n\nHelp the open source community by adding them to\nhttps://lrclib.net/",
		Duration: time.Duration(song.DurationMS) * time.Millisecond,
	}

	return &Missing{song: song, lyrics: lyric}
}

// GetSong returns the song associated to the lyrics
func (m *Missing) GetSong() model.Song {
	return m.song
}

// Previous provides the previous `amount` of lyrics without affecting the current lyric
// In this case it's alway nothing
func (m *Missing) Previous(_ int) []Lyric {
	return []Lyric{}
}

// Current provides the current lyric if any.
func (m *Missing) Current() (Lyric, bool) {
	return m.lyrics, true
}

// Next provides the next lyric.
// In this case it's always nothing
func (m *Missing) Next() (Lyric, bool) {
	return Lyric{}, false
}

// Upcoming provides the next `amount` lyrics without affecting the current lyric
// In this case it's always empty
func (m *Missing) Upcoming(_ int) []Lyric {
	return []Lyric{}
}

// Progress shows the fraction of lyrics that have been used.
func (m *Missing) Progress() float64 {
	return 1
}
