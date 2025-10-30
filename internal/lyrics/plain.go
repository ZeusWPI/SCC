package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

type Plain struct {
	song   model.Song
	lyrics Lyric
}

func newPlain(song model.Song) Lyrics {
	lyric := Lyric{
		Text:     song.Lyrics,
		Duration: time.Duration(song.DurationMS) * time.Millisecond,
	}
	return &Plain{song: song, lyrics: lyric}
}

func (p *Plain) GetSong() model.Song {
	return p.song
}

func (p *Plain) Previous(_ int) []Lyric {
	return []Lyric{}
}

func (p *Plain) Current() (Lyric, bool) {
	return p.lyrics, true
}

func (p *Plain) Next() (Lyric, bool) {
	return Lyric{}, false
}

func (p *Plain) Upcoming(_ int) []Lyric {
	return []Lyric{}
}

func (p *Plain) Progress() float64 {
	return 1
}
