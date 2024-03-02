package api

import (
	"fmt"
	"net/http"
	"slices"

	gin "github.com/gin-gonic/gin"
)

type message struct {
	Message string `form:"message" json:"message" xml:"message" binding:"required"`
}

type header struct {
	Name string `header:"X-Username"`
	Ip   string `header:"X-Real-IP"`
}

var messages uint64 = 0
var blockedNames = []string{"Paul-Henri Spaak"}
var blockedIps = []string{}
var maxMessageLength = 200

func getMessage(c *gin.Context) {
	c.JSON(200, gin.H{"messages": messages})
}

func postMessage(c *gin.Context) {
	// Get structs
	header := &header{}
	message := &message{}

	// Check Header
	if err := c.ShouldBindHeader(header); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check Data
	if err := c.ShouldBindJSON(message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Max message length
	if len(message.Message) > maxMessageLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message too long, maximum " + fmt.Sprint(maxMessageLength)})
		return
	}

	// Check if sender is blocked and construct message
	var newMessage string
	if header.Name != "" {
		if slices.Contains(blockedNames, header.Name) {
			c.JSON(http.StatusOK, gin.H{"message": "Message received"})
			return
		}
		newMessage = fmt.Sprintf("[%s] %s", header.Name, message.Message)
	} else if header.Ip != "" {
		if slices.Contains(blockedIps, header.Ip) {
			c.JSON(http.StatusOK, gin.H{"message": "Message received"})
			return
		}
		newMessage = fmt.Sprintf("<%s> %s", header.Ip, message.Message)
	} else {
		newMessage = message.Message
	}

	// Increment messages
	messages++

	c.JSON(http.StatusOK, gin.H{"message": "Message received"})

	// Add message to cammie chat
	safeMessageQueue.AddMessage(newMessage)
}
