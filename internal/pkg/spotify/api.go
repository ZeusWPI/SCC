package spotify

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/util"
	"go.uber.org/zap"
)

const apiURL = "https://api.spotify.com/v1"

type trackArtist struct {
	Name string `json:"name"`
}

type trackResponse struct {
	Name       string        `json:"name"`
	Artists    []trackArtist `json:"artists"`
	DurationMS int64         `json:"duration_ms"`
}

func (s *Spotify) setTrack(track *dto.Spotify) error {
	zap.S().Info("Spotify: Getting track info for id: ", track.SpotifyID)

	req := fiber.Get(fmt.Sprintf("%s/%s/%s", apiURL, "tracks", track.SpotifyID)).
		Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	res := new(trackResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Spotify: Track request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return errors.New("error getting track")
	}

	track.Title = res.Name
	track.Artists = util.SliceStringJoin(res.Artists, ", ", func(a trackArtist) string { return a.Name })
	track.DurationMS = res.DurationMS

	return nil
}
