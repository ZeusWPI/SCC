// Package song provides all song related logic
package song

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	if s.ClientID == "" || s.ClientSecret == "" {
		return errors.New("Song: Spotify client id or secret not set")
	}

	// Check if song is already in DB
	trackDB, err := s.db.Queries.GetSongBySpotifyID(context.Background(), track.SpotifyID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if (trackDB != sqlc.Song{}) {
		// Already in DB
		// Add to song history if it's not the latest song
		songHistory, err := s.db.Queries.GetLastSongHistory(context.Background())
		if err != nil && err != sql.ErrNoRows {
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

	// Set track info
	err = s.getTrack(track)
	if err != nil {
		return err
	}

	// Store track in DB
	trackDB, err = s.db.Queries.CreateSong(context.Background(), track.CreateParams())
	if err != nil {
		return err
	}

	// Add to song history
	_, err = s.db.Queries.CreateSongHistory(context.Background(), trackDB.ID)

	return err

}
