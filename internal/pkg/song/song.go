// Package song provides all song related logic
package song

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
)

// Song represents a song instance
type Song struct {
	db           *db.DB
	ClientID     string
	ClientSecret string
	AccessToken  string
	ExpiresTime  int64
}

// New creates a new song instance
func New(db *db.DB) (*Song, error) {
	clientID := config.GetDefaultString("song.spotify_client_id", "")
	clientSecret := config.GetDefaultString("song.spotify_client_secret", "")

	if clientID == "" || clientSecret == "" {
		return &Song{}, errors.New("Song: Spotify client id or secret not set")
	}

	return &Song{db: db, ClientID: clientID, ClientSecret: clientSecret, ExpiresTime: 0}, nil
}

// Track gets information about the current track and stores it in the database
func (s *Song) Track(track *dto.Song) error {
	var errs []error

	if s.ClientID == "" || s.ClientSecret == "" {
		return errors.New("Song: Spotify client id or secret not set")
	}

	// Check if song is already in DB
	trackDB, err := s.db.Queries.GetSongBySpotifyID(context.Background(), track.SpotifyID)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	if (trackDB != sqlc.Song{}) {
		// Already in DB
		// Add to song history if it's not the latest song
		songHistory, err := s.db.Queries.GetLastSongHistory(context.Background())
		if err != nil && err != pgx.ErrNoRows {
			return err
		}

		if (songHistory != sqlc.SongHistory{}) && songHistory.SongID == trackDB.ID {
			// Song is already the latest, don't add it again
			return nil
		}

		_, err = s.db.Queries.CreateSongHistory(context.Background(), trackDB.ID)
		return err
	}

	// Not in database yet, add it

	// Refresh token if needed
	if s.ExpiresTime <= time.Now().Unix() {
		err := s.refreshToken()
		if err != nil {
			return err
		}
	}

	// Get track info
	if err = s.getTrack(track); err != nil {
		return err
	}

	// Get lyrics
	if err = s.getLyrics(track); err != nil {
		errs = append(errs, err)
	}

	// Store track in DB
	trackDB, err = s.db.Queries.CreateSong(context.Background(), *track.CreateSongParams())
	if err != nil {
		errs = append(errs, err)
		return errors.Join(errs...)
	}
	track.ID = trackDB.ID

	// Handle artists
	for i, artist := range track.Artists {
		a, err := s.db.Queries.GetSongArtistBySpotifyID(context.Background(), artist.SpotifyID)
		if err != nil && err != pgx.ErrNoRows {
			errs = append(errs, err)
			continue
		}

		if (a != sqlc.SongArtist{}) {
			// Artist already exists
			// Add it as an artist for this track
			if _, err := s.db.Queries.CreateSongArtistSong(context.Background(), *track.CreateSongArtistSongParams(i)); err != nil {
				errs = append(errs, err)
			}
			continue
		}

		// Get artist data
		if err := s.getArtist(&track.Artists[i]); err != nil {
			errs = append(errs, err)
			continue
		}

		// Insert artist in DB
		a, err = s.db.Queries.CreateSongArtist(context.Background(), *track.CreateSongArtistParams(i))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		track.Artists[i].ID = a.ID

		// Add artist as an artist for this song
		if _, err := s.db.Queries.CreateSongArtistSong(context.Background(), *track.CreateSongArtistSongParams(i)); err != nil {
			errs = append(errs, err)
			continue
		}

		// Check if the artists genres are in db
		for j, genre := range track.Artists[i].Genres {
			g, err := s.db.Queries.GetSongGenreByName(context.Background(), genre.Genre)
			if err != nil && err != pgx.ErrNoRows {
				errs = append(errs, err)
				continue
			}

			if (g != sqlc.SongGenre{}) {
				// Genre already exists
				continue
			}

			// Insert genre in DB
			g, err = s.db.Queries.CreateSongGenre(context.Background(), track.CreateSongGenreParams(i, j))
			if err != nil {
				errs = append(errs, err)
				continue
			}
			track.Artists[i].Genres[j].ID = g.ID

			// Add genre as a genre for this artist
			if _, err := s.db.Queries.CreateSongArtistGenre(context.Background(), *track.CreateSongArtistGenreParamas(i, j)); err != nil {
				errs = append(errs, err)
			}
		}
	}

	// Add to song history
	if _, err = s.db.Queries.CreateSongHistory(context.Background(), trackDB.ID); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)

}
