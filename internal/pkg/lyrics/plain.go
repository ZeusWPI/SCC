package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/dto"
)

// Plain represents lyrics that don't have timestamps or songs without lyrics
type Plain struct {
	song   dto.Song
	lyrics Lyric
	given  bool
}

func newPlain(song *dto.Song) Lyrics {
	lyric := Lyric{
		Text:     song.Lyrics,
		Duration: time.Duration(song.DurationMS),
	}
	return &Plain{song: *song, lyrics: lyric, given: false}
}

// GetSong returns the song associated to the lyrics
func (p *Plain) GetSong() dto.Song {
	return p.song
}

// Previous provides the previous `amount` of lyrics without affecting the current lyric
func (p *Plain) Previous(amount int) []Lyric {
	return []Lyric{}
}

// Next provides the next lyric.
// If the lyrics are finished the boolean is set to false
func (p *Plain) Next() (Lyric, bool) {
	if p.given {
		return Lyric{}, false
	}

	return p.lyrics, true
}

// Upcoming provides the next `amount` lyrics without affecting the current lyric
func (p *Plain) Upcoming(amount int) []Lyric {
	return []Lyric{}
}
