package api

import (
	"io"
	"scc/screen"

	"github.com/gin-gonic/gin"
)

func Start(screenApp *screen.ScreenApp) {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	r.GET("/message", getMessage)
	r.POST("/message", postMessage)

	r.Run()
}
