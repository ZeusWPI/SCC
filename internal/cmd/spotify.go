package cmd

import (
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/spotify"
)

// Spotify starts the Spotify integration
func Spotify(db *db.DB) (*spotify.Spotify, error) {
	spotify, err := spotify.New(db)

	return spotify, err
}
