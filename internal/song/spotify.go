package song

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/database/model"
	"go.uber.org/zap"
)

type accountResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (c *client) refreshToken() error {
	zap.S().Info("Song: Refreshing spotify access token")

	form := &fiber.Args{}
	form.Add("grant_type", "client_credentials")

	req := fiber.Post(apiAccount).Form(form).BasicAuth(c.clientID, c.clientSecret)

	res := new(accountResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("spotify token refresh request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return errors.New("getting spotify account access token")
	}

	c.accessToken = res.AccessToken
	c.expiresTime = time.Now().Unix() + res.ExpiresIn

	return nil
}

type trackArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type trackAlbum struct {
	Name string `json:"name"`
}

type trackResponse struct {
	Name       string        `json:"name"`
	Album      trackAlbum    `json:"album"`
	Artists    []trackArtist `json:"artists"`
	DurationMS int32         `json:"duration_ms"`
}

func (c *client) populateSong(song *model.Song) error {
	zap.S().Info("Song: Getting track info for id: ", song.SpotifyID)

	req := fiber.Get(fmt.Sprintf("%s/%s/%s", apiSpotify, "tracks", song.SpotifyID)).
		Set("Authorization", "Bearer "+c.accessToken)

	res := new(trackResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("track request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return fmt.Errorf("track request wrong status code %d", status)
	}

	song.Title = res.Name
	song.Album = res.Album.Name
	song.DurationMS = int(res.DurationMS)

	for _, a := range res.Artists {
		song.Artists = append(song.Artists, model.Artist{
			Name:      a.Name,
			SpotifyID: a.ID,
		})
	}

	return nil
}

type artistResponse struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

func (c *client) populateArtist(artist *model.Artist) error {
	zap.S().Info("Song: Getting artists info for ", artist.SpotifyID)

	req := fiber.Get(fmt.Sprintf("%s/%s/%s", apiSpotify, "artists", artist.SpotifyID)).
		Set("Authorization", "Bearer "+c.accessToken)

	res := new(artistResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("artist request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return fmt.Errorf("artist request wrong status code %d", status)
	}

	for _, genre := range res.Genres {
		artist.Genres = append(artist.Genres, model.Genre{Genre: genre})
	}

	return nil
}
