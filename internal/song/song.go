package song

import (
	"errors"
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

const (
	apiSpotify = "https://api.spotify.com/v1"
	apiAccount = "https://accounts.spotify.com/api/token"
	apiLrc     = "https://lrclib.net/api"
)

type client struct {
	clientID     string
	clientSecret string

	accessToken string
	expiresTime int64
}

var C *client

func Init() error {
	clientID := config.GetDefaultString("backend.song.spotify_client_id", "")
	clientSecret := config.GetDefaultString("backend.song.spotify_client_secret", "")

	zap.S().Info(clientID)

	if clientID == "" || clientSecret == "" {
		return errors.New("spotify client id or secret not set")
	}

	C = &client{
		clientID:     clientID,
		clientSecret: clientSecret,
		expiresTime:  0,
	}

	return nil
}

func (c *client) Populate(song *model.Song) error {
	zap.S().Info("Populating song")

	if c.expiresTime <= time.Now().Unix() {
		err := c.refreshToken()
		if err != nil {
			return err
		}
	}

	if err := c.populateSong(song); err != nil {
		return err
	}

	for i := range song.Artists {
		if err := c.populateArtist(&song.Artists[i]); err != nil {
			return err
		}
	}

	if err := c.getLyrics(song); err != nil {
		return err
	}

	zap.S().Info("Populated song")

	return nil
}
