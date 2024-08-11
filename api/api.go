package api

import (
	"io"
	"scc/screen"

	"github.com/gin-gonic/gin"
)

// Wrapper for the handler functions to pass the screen application
func handlerWrapper(app *screen.ScreenApp, callback func(*screen.ScreenApp, *gin.Context)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		callback(app, ctx)
	}
}

// Start the API
func Start(screenApp *screen.ScreenApp) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	// Routes

	// Cammie chat routes
	r.GET("/message", handlerWrapper(screenApp, cammieGetMessage))
	r.POST("/message", handlerWrapper(screenApp, cammiePostMessage))

	// Spotify routes
	r.POST("/spotify", handlerWrapper(screenApp, spotifyGetMessage))

	// Start Tap
	go tapRunRequests(screenApp)

	// Start Zess
	go zessRunRequests(screenApp)

	// Start API
	r.Run()
}
