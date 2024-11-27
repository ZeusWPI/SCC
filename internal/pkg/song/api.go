package song

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/util"
	"go.uber.org/zap"
)

type trackArtist struct {
	Name string `json:"name"`
}

type trackResponse struct {
	Name       string        `json:"name"`
	Artists    []trackArtist `json:"artists"`
	DurationMS int64         `json:"duration_ms"`
}

func (s *Song) getTrack(track *dto.Song) error {
	zap.S().Info("Song: Getting track info for id: ", track.SpotifyID)

	api := config.GetDefaultString("song.spotify_api", "https://api.spotify.com/v1")
	req := fiber.Get(fmt.Sprintf("%s/%s/%s", api, "tracks", track.SpotifyID)).
		Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	res := new(trackResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Song: Track request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return errors.New("Song: Error getting track")
	}

	track.Title = res.Name
	track.Artists = util.SliceStringJoin(res.Artists, ", ", func(a trackArtist) string { return a.Name })
	track.DurationMS = res.DurationMS

	return nil
}
