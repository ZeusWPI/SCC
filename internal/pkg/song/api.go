package song

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"go.uber.org/zap"
)

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

func (s *Song) getTrack(track *dto.Song) error {
	zap.S().Info("Song: Getting track info for id: ", track.SpotifyID)

	req := fiber.Get(fmt.Sprintf("%s/%s/%s", s.api, "tracks", track.SpotifyID)).
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
	track.Album = res.Album.Name
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
	zap.S().Info("Song: Getting artists info for ", artist.ID)

	req := fiber.Get(fmt.Sprintf("%s/%s/%s", s.api, "artists", artist.SpotifyID)).
		Set("Authorization", fmt.Sprintf("Bearer %s", s.AccessToken))

	res := new(artistResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Song: Artist request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return fmt.Errorf("Song: Artist request wrong status code %d", status)
	}

	artist.Popularity = int32(res.Popularity)
	artist.Followers = int32(res.Followers.Total)

	for _, genre := range res.Genres {
		artist.Genres = append(artist.Genres, dto.SongGenre{Genre: genre})
	}

	return nil
}

type lyricsResponse struct {
	PlainLyrics  string `json:"plainLyrics"`
	SyncedLyrics string `json:"SyncedLyrics"`
}

func (s *Song) getLyrics(track *dto.Song) error {
	zap.S().Info("Song: Getting lyrics for ", track.Title)

	// Get most popular artist
	if len(track.Artists) == 0 {
		return fmt.Errorf("Song: No artists for track: %v", track)
	}
	artist := track.Artists[0]
	for _, a := range track.Artists {
		if a.Followers > artist.Followers {
			artist = a
		}
	}

	// Construct url
	params := url.Values{}
	params.Set("artist_name", artist.Name)
	params.Set("track_name", track.Title)
	params.Set("album_name", track.Album)
	params.Set("duration", fmt.Sprintf("%d", track.DurationMS/1000))

	req := fiber.Get(fmt.Sprintf("%s/get?%s", s.apiLrclib, params.Encode()))

	res := new(lyricsResponse)
	status, _, errs := req.Struct(res)
	if len(errs) > 0 {
		return errors.Join(append([]error{errors.New("Song: Lyrics request failed")}, errs...)...)
	}
	if status != fiber.StatusOK {
		return fmt.Errorf("Song: Lyrics request wrong status code %d", status)
	}
	if (res == &lyricsResponse{}) {
		return errors.New("Song: Lyrics request returned empty struct")
	}

	if res.SyncedLyrics != "" {
		track.LyricsType = "synced"
		track.Lyrics = res.SyncedLyrics
	} else {
		track.LyricsType = "plain"
		track.Lyrics = res.PlainLyrics
	}

	return nil
}
