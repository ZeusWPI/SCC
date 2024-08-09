package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
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

var accessToken = ""
var expiresOn int64 = 0
var clientID = "d385173507a54bca93cc3327c0c2f5d9"
var clientSecret = "8e78977c1ba54b90b17f9dcd6b301c37"

func spotifyHandlerWrapper(app *screen.ScreenApp) func(*gin.Context) {
	return func(ctx *gin.Context) {
		spotifyHandler(app, ctx)
	}
}

func spotifyHandler(app *screen.ScreenApp, ctx *gin.Context) {
	message := &spotifyMessage{}

	if err := ctx.ShouldBindJSON(message); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"track_id": "Track ID received"})

	if expiresOn < time.Now().Unix() {
		if err := setAccessToken(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to refresh spotify token: %s", err)
		}
	}

	track, err := getTrackTitle(message.TrackID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to get track information: %s", err)
	}

	app.Spotify.Update(track)
}

func setAccessToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

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

	accessToken = message.AccessToken
	expiresOn = time.Now().Unix() + message.ExpiresIn

	return nil
}

func getTrackTitle(trackID string) (string, error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/tracks/%s", trackID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	trackResponse := &spotifyTrackResponse{}
	if err := json.NewDecoder(resp.Body).Decode(trackResponse); err != nil {
		return "", err
	}

	trackTitle := trackResponse.Name
	artistsNames := make([]string, len(trackResponse.Artists))
	for i, artist := range trackResponse.Artists {
		artistsNames[i] = artist.Name
	}

	return fmt.Sprintf("%s - %s", trackTitle, strings.Join(artistsNames, ", ")), nil
}
