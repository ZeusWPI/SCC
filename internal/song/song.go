// Package song provides all song related logic
package song

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
)

type Song struct {
	song repository.Song

	clientID     string
	clientSecret string
	accessToken  string
	expiresTime  int64

	url        string
	urlAccount string
	urlLrclib  string
}

func New(repo repository.Repository) (*Song, error) {
	clientID := config.GetDefaultString("backend.song.spotify_client_id", "")
	clientSecret := config.GetDefaultString("backend.song.spotify_client_secret", "")

	if clientID == "" || clientSecret == "" {
		return &Song{}, errors.New("song: spotify client id or secret not set")
	}

	return &Song{
		song:         *repo.NewSong(),
		clientID:     clientID,
		clientSecret: clientSecret,
		expiresTime:  0,
		url:          config.GetDefaultString("backend.song.spotify_url", "https://api.spotify.com/v1"),
		urlAccount:   config.GetDefaultString("backend.song.spotify_url_account", "https://accounts.spotify.com/api/token"),
		urlLrclib:    config.GetDefaultString("backend.song.lrclib_url", "https://lrclib.net/api"),
	}, nil
}

// Track gets information about the current track and stores it in the database
func (s *Song) Track(ctx context.Context, track *dto.Song) error {
	var errs []error

	if s.clientID == "" || s.clientSecret == "" {
		return errors.New("Song: Spotify client id or secret not set")
	}

	// Check if song is already in DB
	trackDB, err := s.song.GetBySpotifyID(ctx, track.SpotifyID)
	if err != nil {
		return err
	}

	if trackDB != nil {
		// Already in DB
		// Add to song history if it's not the latest song
		songHistory, err := s.song.GetLast(ctx)
		if err != nil {
			return err
		}

		if songHistory != nil && songHistory.ID == trackDB.ID {
			// Song is already the latest, don't add it again
			return nil
		}

		if err = s.song.CreateHistory(ctx, trackDB.ID); err != nil {
			return err
		}

		return nil
	}

	// Not in database yet, add it

	// Refresh token if needed
	if s.expiresTime <= time.Now().Unix() {
		err := s.refreshToken()
		if err != nil {
			return err
		}
	}

	// Get track info
	if err = s.getTrack(ctx, track); err != nil {
		return err
	}

	// Get lyrics
	if err = s.getLyrics(ctx, track); err != nil {
		errs = append(errs, err)
	}

	// Store track in DB
	err = s.song.Create(ctx, &track)
	if err != nil {
		errs = append(errs, err)
		return errors.Join(errs...)
	}
	track.ID = trackDB.ID

	// Handle artists
	for i, artist := range track.Artists {
		a, err := s.song.GetArtistBySpotifyID(ctx, artist.SpotifyID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if a != nil {
			// Artist already exists
			// Add it as an artist for this track
			track.Artists[i].ID = a.ID
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
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
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
