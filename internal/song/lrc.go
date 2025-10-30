package song

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/database/model"
	"go.uber.org/zap"
)

type lyricsResponse struct {
	Instrumental bool   `json:"instrumental"`
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"SyncedLyrics"`
}

func (c *client) getLyrics(song *model.Song) error {
	zap.S().Info("Song: Getting lyrics for ", song.Title)

	// Get most popular artist
	if len(song.Artists) == 0 {
		return fmt.Errorf("no artists for track: %v", song)
	}
	artist := song.Artists[0]

	// Construct url
	params := url.Values{}
	params.Set("artist_name", artist.Name)
	params.Set("track_name", song.Title)
	params.Set("album_name", song.Album)
	params.Set("duration", strconv.Itoa(song.DurationMS/1000))

	req := fiber.Get(fmt.Sprintf("%s/get?%s", apiLrc, params.Encode()))

	res := new(lyricsResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("lyrics request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		if status == fiber.StatusNotFound {
			// Lyrics not found
			song.LyricsType = model.LyricsMissing
			return nil
		}

		return fmt.Errorf("lyrics request wrong status code %d", status)
	}
	if (res == &lyricsResponse{}) {
		return errors.New("lyrics request returned empty struct")
	}

	switch {
	case res.SyncedLyrics != "":
		song.LyricsType = model.LyricsSynced
		song.Lyrics = res.SyncedLyrics
	case res.PlainLyrics != "":
		song.LyricsType = model.LyricsPlain
		song.Lyrics = res.PlainLyrics
	case res.Instrumental:
		song.LyricsType = model.LyricsInstrumental
		song.Lyrics = ""
	default:
		song.LyricsType = model.LyricsMissing
	}

	return nil
}
