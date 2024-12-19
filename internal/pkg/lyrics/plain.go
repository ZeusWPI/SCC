package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/dto"
)

// Plain represents lyrics that don't have timestamps or songs without lyrics
type Plain struct {
	song   dto.Song
	lyrics Lyric
}

func newPlain(song dto.Song) Lyrics {
	lyric := Lyric{
		Text:     song.Lyrics,
		Duration: time.Duration(song.DurationMS) * time.Millisecond,
	}
	return &Plain{song: song, lyrics: lyric}
}

// GetSong returns the song associated to the lyrics
func (p *Plain) GetSong() dto.Song {
	return p.song
}

// Previous provides the previous `amount` of lyrics without affecting the current lyric
// In this case it's always nothing
func (p *Plain) Previous(_ int) []Lyric {
	return []Lyric{}
}

// Current provides the current lyric if any.
func (p *Plain) Current() (Lyric, bool) {
	return p.lyrics, true
}

// Next provides the next lyric.
// In this case it's alway nothing
func (p *Plain) Next() (Lyric, bool) {
	return Lyric{}, false
}

// Upcoming provides the next `amount` lyrics without affecting the current lyric
// In this case it's always empty
func (p *Plain) Upcoming(_ int) []Lyric {
	return []Lyric{}
}

// Progress shows the fraction of lyrics that have been used.
func (p *Plain) Progress() float64 {
	return 1
}
