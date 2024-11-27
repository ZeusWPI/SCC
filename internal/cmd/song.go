package cmd

import (
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/song"
)

// Song starts the Song integration
func Song(db *db.DB) (*song.Song, error) {
	song, err := song.New(db)

	return song, err
}
