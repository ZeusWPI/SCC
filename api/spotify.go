package api

import (
	"github.com/gin-gonic/gin"
	"scc/screen"
)

func spotifyHandlerWrapper(app *screen.ScreenApp) func(*gin.Context) {
	return func(ctx *gin.Context) {
		spotifyHandler(app, ctx)
	}
}

func spotifyHandler(app *screen.ScreenApp, ctx *gin.Context) {
	b, _ := ctx.GetRawData()
	app.Spotify.Update(string(b))
}
