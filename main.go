package main

import (
	"scc/api"

	gin "github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/message", api.GetMessage)
	r.POST("/message", api.PostMessage)

	r.Run()
}
