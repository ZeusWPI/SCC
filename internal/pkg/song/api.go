package song

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

var api = config.GetDefaultString("song.spotify_api", "https://api.spotify.com/v1")

type trackArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type trackResponse struct {
	Name       string        `json:"name"`
	Artists    []trackArtist `json:"artists"`
	DurationMS int64         `json:"duration_ms"`
}

func (s *Song) getTrack(track *dto.Song) error {
	zap.S().Info("Song: Getting track info for id: ", track.SpotifyID)

	req := fiber.Get(fmt.Sprintf("%s/%s/%s", api, "tracks", track.SpotifyID)).
		Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	res := new(trackResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Song: Track request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return fmt.Errorf("Song: Track request wrong status code %d", status)
	}

	track.Title = res.Name
	track.DurationMS = res.DurationMS

	for _, a := range res.Artists {
		track.Artists = append(track.Artists, dto.SongArtist{
			Name:      a.Name,
			SpotifyID: a.ID,
		})
	}

	return nil
}

type artistFollowers struct {
	Total int `json:"total"`
}

type artistResponse struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Genres     []string        `json:"genres"`
	Popularity int             `json:"popularity"`
	Followers  artistFollowers `json:"followers"`
}

func (s *Song) getArtist(artist *dto.SongArtist) error {
	req := fiber.Get(fmt.Sprintf("%s/%s/%s", api, "artists", artist.SpotifyID)).
		Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	res := new(artistResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Song: Artist request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return fmt.Errorf("Song: Artist request wrong status code %d", status)
	}

	artist.Popularity = int64(res.Popularity)
	artist.Followers = int64(res.Followers.Total)

	for _, genre := range res.Genres {
		artist.Genres = append(artist.Genres, dto.SongGenre{Genre: genre})
	}

	return nil
}
