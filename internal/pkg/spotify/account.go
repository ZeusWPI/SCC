package spotify

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const accountURL = "https://accounts.spotify.com/api/token"

type accountResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (s *Spotify) refreshToken() error {
	zap.S().Info("Spotify: Refreshing access token")

	body, err := json.Marshal(fiber.Map{
		"grant_type":    "client_credentials",
		"client_id":     s.ClientID,
		"client_secret": s.ClientSecret,
	})
	if err != nil {
		return err
	}

	req := fiber.Post(accountURL).Body(body).ContentType("application/json")

	res := new(accountResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	if status != fiber.StatusOK {
		return errors.New("error getting access token")
	}

	s.AccessToken = res.AccessToken
	s.ExpiresTime = time.Now().Unix() + res.ExpiresIn

	return nil
}
