// Package spotify provides all spotify related logic
package spotify

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

// Spotify represents a spotify instance
type Spotify struct {
	db           *db.DB
	ClientID     string
	ClientSecret string
	AccessToken  string
	ExpiresTime  int64
}

// New creates a new spotify instance
func New(db *db.DB) (*Spotify, error) {
	clientID := config.GetDefaultString("spotify.client_id", "")
	clientSecret := config.GetDefaultString("spotify.client_secret", "")

	if clientID == "" || clientSecret == "" {
		return &Spotify{}, errors.New("Spotify client id or secret not set")
	}

	return &Spotify{db: db, ClientID: clientID, ClientSecret: clientSecret, ExpiresTime: 0}, nil
}

// Track gets information about the current track and stores it in the database
func (s *Spotify) Track(track *dto.Spotify) error {
	if s.ClientID == "" || s.ClientSecret == "" {
		return errors.New("spotify client id or secret not set")
	}

	// Check if song is already in DB
	trackDB, err := s.db.Queries.GetSpotifyBySpotifyID(context.Background(), track.SpotifyID)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if (trackDB != sqlc.Spotify{}) {
		// Already in DB
		// No need to refetch data
		track.Title = trackDB.Title
		track.Artists = trackDB.Artists
		track.DurationMS = trackDB.DurationMs
		_, err := s.db.Queries.CreateSpotify(context.Background(), track.CreateParams())

		return err
	}

	// Refresh token if needed
	if s.ExpiresTime <= time.Now().Unix() {
		err := s.refreshToken()
		if err != nil {
			return err
		}
	}

	// Set track info
	err = s.setTrack(track)
	if err != nil {
		return err
	}

	// Store track in DB
	_, err = s.db.Queries.CreateSpotify(context.Background(), track.CreateParams())
	if err != nil {
		return err
	}

	return nil

}
