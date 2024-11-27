package dto

import (
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Song is the DTO for the song
type Song struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Artists    string `json:"artists"`
	SpotifyID  string `json:"spotify_id" validate:"required"`
	DurationMS int64  `json:"duration_ms"`
}

// SongDTO converts a sqlc.Song to a Song
func SongDTO(song sqlc.Song) *Song {
	return &Song{
		ID:         song.ID,
		Title:      song.Title,
		Artists:    song.Artists,
		SpotifyID:  song.SpotifyID,
		DurationMS: song.DurationMs,
	}
}

// CreateParams converts a Song to sqlc.CreateSongParams
func (s *Song) CreateParams() sqlc.CreateSongParams {
	return sqlc.CreateSongParams{
		Title:      s.Title,
		Artists:    s.Artists,
		SpotifyID:  s.SpotifyID,
		DurationMs: s.DurationMS,
	}
}
