package dto

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Spotify is the DTO for the spotify
type Spotify struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Artists    string    `json:"artists"`
	SpotifyID  string    `json:"spotify_id" validate:"required"`
	DurationMS int64     `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

// SpotifyDTO converts a sqlc.Spotify to a Spotify
func SpotifyDTO(spotify sqlc.Spotify) *Spotify {
	return &Spotify{
		ID:         spotify.ID,
		Title:      spotify.Title,
		Artists:    spotify.Artists,
		SpotifyID:  spotify.SpotifyID,
		DurationMS: spotify.DurationMs,
		CreatedAt:  spotify.CreatedAt,
	}
}

// CreateParams converts a Spotify to sqlc.CreateSpotifyParams
func (s *Spotify) CreateParams() sqlc.CreateSpotifyParams {
	return sqlc.CreateSpotifyParams{
		Title:      s.Title,
		Artists:    s.Artists,
		SpotifyID:  s.SpotifyID,
		DurationMs: s.DurationMS,
	}
}
