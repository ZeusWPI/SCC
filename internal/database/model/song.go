package model

import (
	"time"

	"github.com/zeusWPI/scc/pkg/sqlc"
)

type SongGenre struct {
	ID    int
	Genre string
}

type SongArtist struct {
	ID         int
	Name       string
	SpotifyID  string
	Followers  int
	Popularity int
	Genres     []SongGenre
}

type Song struct {
	ID         int
	Title      string
	Album      string
	SpotifyID  string
	DurationMS int
	LyricsType string
	Lyrics     string
	CreatedAt  time.Time
	Artists    []SongArtist
}

func SongModel(s sqlc.Song) *Song {
	var lyricsType string
	if s.LyricsType.Valid {
		lyricsType = s.Lyrics.String
	}

	var lyrics string
	if s.Lyrics.Valid {
		lyrics = s.Lyrics.String
	}

	return &Song{
		ID:         int(s.ID),
		Title:      s.Title,
		Album:      s.Album,
		SpotifyID:  s.SpotifyID,
		DurationMS: int(s.DurationMs),
		LyricsType: lyricsType,
		Lyrics:     lyrics,
	}
}

func SongModelHistory(songs []sqlc.GetLastSongFullRow) *Song {
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
				ID:         int(song.ArtistID.Int32),
				Name:       song.ArtistName.String,
				SpotifyID:  song.ArtistSpotifyID.String,
				Followers:  int(song.ArtistFollowers.Int32),
				Popularity: int(song.ArtistPopularity.Int32),
				Genres:     make([]SongGenre, 0),
			}
			artistsMap[song.ArtistID.Int32] = artist
		}

		// Add genre
		artist.Genres = append(artist.Genres, SongGenre{
			ID:    int(song.GenreID.Int32),
			Genre: song.Genre.String,
		})
	}

	artists := make([]SongArtist, 0, len(artistsMap))
	for _, artist := range artistsMap {
		artists = append(artists, artist)
	}

	return &Song{
		ID:         int(songs[0].ID),
		Title:      songs[0].SongTitle,
		Album:      songs[0].Album,
		SpotifyID:  songs[0].SpotifyID,
		DurationMS: int(songs[0].DurationMs),
		LyricsType: lyricsType,
		Lyrics:     lyrics,
		CreatedAt:  songs[0].CreatedAt.Time,
		Artists:    artists,
	}
}
