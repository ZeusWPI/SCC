package api

import (
	"io"

	"scc/screen"

	"github.com/gin-gonic/gin"
)

var safeMessageQueue *screen.SafeMessageQueue

func Start(queue *screen.SafeMessageQueue) {
	safeMessageQueue = queue

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	r.GET("/message", getMessage)
	r.POST("/message", postMessage)

	r.Run()
}
