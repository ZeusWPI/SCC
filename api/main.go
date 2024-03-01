package api

import (
	"io"

	"github.com/gin-gonic/gin"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()

	r.GET("/message", getMessage)
	r.POST("/message", postMessage)

	r.Run()
}
