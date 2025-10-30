package model

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/sqlc"
)

type LyricsType string

const (
	LyricsPlain        LyricsType = "plain"
	LyricsSynced       LyricsType = "synced"
	LyricsInstrumental LyricsType = "instrumental"
	LyricsMissing      LyricsType = "missing"
)

type Song struct {
	ID         int
	Title      string
	Album      string
	SpotifyID  string
	DurationMS int
	LyricsType LyricsType
	Lyrics     string
	PlayedAt   time.Time
	Artists    []Artist
}

type Artist struct {
	ID        int
	Name      string
	SpotifyID string
	Genres    []Genre
}

type Genre struct {
	ID    int
	Genre string
}

func SongModel(s sqlc.Song) *Song {
	lyrics := ""
	if s.Lyrics.Valid {
		lyrics = s.Lyrics.String
	}

	return &Song{
		ID:         int(s.ID),
		Title:      s.Title,
		Album:      s.Album,
		SpotifyID:  s.SpotifyID,
		DurationMS: int(s.DurationMs),
		LyricsType: LyricsType(s.LyricsType),
		Lyrics:     lyrics,
	}
}

func ArtistModel(a sqlc.SongArtist) *Artist {
	return &Artist{
		ID:        int(a.ID),
		Name:      a.Name,
		SpotifyID: a.SpotifyID,
	}
}

func GenreModel(g sqlc.SongGenre) *Genre {
	return &Genre{
		ID:    int(g.ID),
		Genre: g.Genre,
	}
}
