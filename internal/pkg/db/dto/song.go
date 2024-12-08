package dto

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Song is the DTO for a song
type Song struct {
	ID         int32        `json:"id"`
	Title      string       `json:"title"`
	Album      string       `json:"album"`
	SpotifyID  string       `json:"spotify_id" validate:"required"`
	DurationMS int32        `json:"duration_ms"`
	LyricsType string       `json:"lyrics_type"` // Either 'synced' or 'plain'
	Lyrics     string       `json:"lyrics"`
	CreatedAt  time.Time    `json:"created_at"`
	Artists    []SongArtist `json:"artists"`
}

// SongArtist is the DTO for a song artist
type SongArtist struct {
	ID         int32       `json:"id"`
	Name       string      `json:"name"`
	SpotifyID  string      `json:"spotify_id"`
	Followers  int32       `json:"followers"`
	Popularity int32       `json:"popularity"`
	Genres     []SongGenre `json:"genres"`
}

// SongGenre is the DTO for a song genre
type SongGenre struct {
	ID    int32  `json:"id"`
	Genre string `json:"genre"`
}

// SongDTO converts a sqlc.Song to a Song
func SongDTO(song sqlc.Song) *Song {
	var lyricsType string
	if song.LyricsType.Valid {
		lyricsType = song.Lyrics.String
	}
	var lyrics string
	if song.Lyrics.Valid {
		lyrics = song.Lyrics.String
	}

	return &Song{
		ID:         song.ID,
		Title:      song.Title,
		Album:      song.Album,
		SpotifyID:  song.SpotifyID,
		DurationMS: song.DurationMs,
		LyricsType: lyricsType,
		Lyrics:     lyrics,
	}
}

// SongDTOHistory converts a sqlc.GetLastSongFullRow array to a Song
func SongDTOHistory(songs []sqlc.GetLastSongFullRow) *Song {
	if len(songs) == 0 {
		return nil
	}

	var lyricsType string
	if songs[0].LyricsType.Valid {
		lyricsType = songs[0].LyricsType.String
	}
	var lyrics string
	if songs[0].Lyrics.Valid {
		lyrics = songs[0].Lyrics.String
	}

	artistsMap := make(map[int32]SongArtist)
	for _, song := range songs {
		if !song.ArtistID.Valid {
			continue
		}

		// Get artist
		artist, ok := artistsMap[song.ArtistID.Int32]
		if !ok {
			// Artist doesn't exist yet, add him
			artist = SongArtist{
				ID:         song.ArtistID.Int32,
				Name:       song.ArtistName.String,
				SpotifyID:  song.ArtistSpotifyID.String,
				Followers:  song.ArtistFollowers.Int32,
				Popularity: song.ArtistPopularity.Int32,
				Genres:     make([]SongGenre, 0),
			}
			artistsMap[song.ArtistID.Int32] = artist
		}

		// Add genre
		artist.Genres = append(artist.Genres, SongGenre{
			ID:    song.GenreID.Int32,
			Genre: song.Genre.String,
		})
	}

	artists := make([]SongArtist, 0, len(artistsMap))
	for _, artist := range artistsMap {
		artists = append(artists, artist)
	}

	return &Song{
		ID:         songs[0].ID,
		Title:      songs[0].SongTitle,
		Album:      songs[0].Album,
		SpotifyID:  songs[0].SpotifyID,
		DurationMS: songs[0].DurationMs,
		LyricsType: lyricsType,
		Lyrics:     lyrics,
		CreatedAt:  songs[0].CreatedAt.Time,
		Artists:    artists,
	}
}

// CreateSongParams converts a Song DTO to a sqlc CreateSongParams object
func (s *Song) CreateSongParams() *sqlc.CreateSongParams {
	return &sqlc.CreateSongParams{
		Title:      s.Title,
		Album:      s.Album,
		SpotifyID:  s.SpotifyID,
		DurationMs: s.DurationMS,
		LyricsType: pgtype.Text{String: s.LyricsType, Valid: true},
		Lyrics:     pgtype.Text{String: s.Lyrics, Valid: true},
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
