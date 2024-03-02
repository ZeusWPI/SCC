package api

import (
	"io"
	"scc/screen"

	"github.com/gin-gonic/gin"
)

func handlerWrapper(app *screen.ScreenApp, callback func(*screen.ScreenApp, *gin.Context)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		callback(app, ctx)
	}
}

func Start(screenApp *screen.ScreenApp) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	r.GET("/message", handlerWrapper(screenApp, getMessage))
	r.POST("/message", handlerWrapper(screenApp, postMessage))

	r.Run()
}
