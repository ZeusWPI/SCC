package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"scc/config"
	"scc/screen"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type spotifyMessage struct {
	TrackID string `json:"track_id"`
}

type spotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type spotifyArtist struct {
	Name string `json:"name"`
}

type spotifyTrackResponse struct {
	Name    string          `json:"name"`
	Artists []spotifyArtist `json:"artists"`
}

var (
	spotifyAccessToken        = ""
	spotifyExpiresOn    int64 = 0
	spotifyClientID           = config.GetConfig().Spotify.ClientID
	spotifyClientSecret       = config.GetConfig().Spotify.ClientSecret
)

func spotifyGetMessage(app *screen.ScreenApp, ctx *gin.Context) {
	message := &spotifyMessage{}

	if err := ctx.ShouldBindJSON(message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"track_id": "Track ID received"})

	if spotifyExpiresOn < time.Now().Unix() {
		if err := spotifySetAccessToken(); err != nil {
			log.Printf("Error: Unable to refresh spotify token: %s\n", err)
		}
	}

	track, err := spotifyGetTrackTitle(message.TrackID)

	if err != nil {
		log.Printf("Error: Unable to get track information: %s\n", err)
	}

	app.Spotify.Update(track)
}

func spotifySetAccessToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", spotifyClientID)
	data.Set("client_secret", spotifyClientSecret)

	// Send the POST request
	resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: received non-200 status code %d", resp.StatusCode)
	}

	message := &spotifyTokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(message); err != nil {
		return err
	}

	spotifyAccessToken = message.AccessToken
	spotifyExpiresOn = time.Now().Unix() + message.ExpiresIn

	return nil
}

func spotifyGetTrackTitle(trackID string) (string, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)
	headers := []header{
		jsonHeader,
		{
			"Authorization",
			"Bearer " + spotifyAccessToken,
		},
	}
	trackResponse := &spotifyTrackResponse{}

	if err := makeGetRequest(url, headers, trackResponse); err != nil {
		return "", err
	}

	trackTitle := trackResponse.Name
	artistsNames := make([]string, len(trackResponse.Artists))
	for i, artist := range trackResponse.Artists {
		artistsNames[i] = artist.Name
	}

	return fmt.Sprintf("%s - %s", trackTitle, strings.Join(artistsNames, ", ")), nil
}
