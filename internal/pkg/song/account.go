package song

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type accountResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (s *Song) refreshToken() error {
	zap.S().Info("Song: Refreshing spotify access token")

	form := &fiber.Args{}
	form.Add("grant_type", "client_credentials")

	api := config.GetDefaultString("backend.song.spotify_api_account", "https://accounts.spotify.com/api/token")
	req := fiber.Post(api).Form(form).BasicAuth(s.ClientID, s.ClientSecret)

	res := new(accountResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Song: Spotify token refresh request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return errors.New("Song: Error getting access token")
	}

	s.AccessToken = res.AccessToken
	s.ExpiresTime = time.Now().Unix() + res.ExpiresIn

	return nil
}
