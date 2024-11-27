package dto

import (
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Song is the DTO for a song
type Song struct {
	ID         int64        `json:"id"`
	Title      string       `json:"title"`
	SpotifyID  string       `json:"spotify_id" validate:"required"`
	DurationMS int64        `json:"duration_ms"`
	Artists    []SongArtist `json:"artists"`
}

// SongArtist is the DTO for a song artist
type SongArtist struct {
	ID         int64       `json:"id"`
	Name       string      `json:"name"`
	SpotifyID  string      `json:"spotify_id"`
	Followers  int64       `json:"followers"`
	Popularity int64       `json:"popularity"`
	Genres     []SongGenre `json:"genres"`
}

// SongGenre is the DTO for a song genre
type SongGenre struct {
	ID    int64  `json:"id"`
	Genre string `json:"genre"`
}

// SongDTO converts a sqlc.Song to a Song
func SongDTO(song sqlc.Song) *Song {
	return &Song{
		ID:         song.ID,
		Title:      song.Title,
		SpotifyID:  song.SpotifyID,
		DurationMS: song.DurationMs,
	}
}

// CreateSongParams converts a Song DTO to a sqlc CreateSongParams object
func (s *Song) CreateSongParams() *sqlc.CreateSongParams {
	return &sqlc.CreateSongParams{
		Title:      s.Title,
		SpotifyID:  s.SpotifyID,
		DurationMs: s.DurationMS,
	}
}

// CreateSongGenreParams converts a Song DTO to a string to create a new genre
func (s *Song) CreateSongGenreParams(idxArtist, idxGenre int) string {
	return s.Artists[idxArtist].Genres[idxGenre].Genre
}

// CreateSongArtistParams converts a Song DTO to a sqlc CreateSongArtistParams object
func (s *Song) CreateSongArtistParams(idxArtist int) *sqlc.CreateSongArtistParams {
	return &sqlc.CreateSongArtistParams{
		Name:       s.Artists[idxArtist].Name,
		SpotifyID:  s.Artists[idxArtist].SpotifyID,
		Followers:  s.Artists[idxArtist].Followers,
		Popularity: s.Artists[idxArtist].Popularity,
	}
}

// CreateSongArtistSongParams converts a Song DTO to a sqlc CreateSongArtistSongParams object
func (s *Song) CreateSongArtistSongParams(idxArtist int) *sqlc.CreateSongArtistSongParams {
	return &sqlc.CreateSongArtistSongParams{
		ArtistID: s.Artists[idxArtist].ID,
		SongID:   s.ID,
	}
}

// CreateSongArtistGenreParamas converts a Song DTO to a sqlc CreateSongArtistGenreParams object
func (s *Song) CreateSongArtistGenreParamas(idxArtist, idxGenre int) *sqlc.CreateSongArtistGenreParams {
	return &sqlc.CreateSongArtistGenreParams{
		ArtistID: s.Artists[idxArtist].ID,
		GenreID:  s.Artists[idxArtist].Genres[idxGenre].ID,
	}
}
