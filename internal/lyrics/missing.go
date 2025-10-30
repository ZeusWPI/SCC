package lyrics

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

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

func (m *Missing) GetSong() model.Song {
	return m.song
}

func (m *Missing) Previous(_ int) []Lyric {
	return []Lyric{}
}

func (m *Missing) Current() (Lyric, bool) {
	return m.lyrics, true
}

func (m *Missing) Next() (Lyric, bool) {
	return Lyric{}, false
}

func (m *Missing) Upcoming(_ int) []Lyric {
	return []Lyric{}
}

func (m *Missing) Progress() float64 {
	return 1
}
