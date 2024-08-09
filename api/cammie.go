package api

import (
	"fmt"
	"net/http"
	"scc/config"
	"scc/screen"
	"slices"

	gin "github.com/gin-gonic/gin"
)

// cammieMessage struct
type cammieMessage struct {
	Message string `form:"message" json:"message" xml:"message" binding:"required"`
}

// cammieHeader struct
type cammieHeader struct {
	Name string `header:"X-Username"`
	IP   string `header:"X-Real-IP"`
}

var (
	cammieMessages         uint64 = 0
	cammieBlockedNames            = config.GetConfig().Cammie.BlockedNames     // Blocked names
	cammieBlockedIps              = config.GetConfig().Cammie.BlockedIps       // Blocked IPs
	cammieMaxMessageLength        = config.GetConfig().Cammie.MaxMessageLength // Maximum message length
)

func cammieGetMessage(app *screen.ScreenApp, c *gin.Context) {
	c.JSON(200, gin.H{"messages": cammieMessages})
}

func cammiePostMessage(app *screen.ScreenApp, c *gin.Context) {
	// Get structs
	header := &cammieHeader{}
	message := &cammieMessage{}

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
	if len(message.Message) > cammieMaxMessageLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message too long, maximum " + fmt.Sprint(cammieMaxMessageLength)})
		return
	}

	// Check if sender is blocked and construct message
	var newMessage string
	if header.Name != "" {
		if slices.Contains(cammieBlockedNames, header.Name) {
			c.JSON(http.StatusOK, gin.H{"message": "Message received"})
			return
		}
		newMessage = fmt.Sprintf("[%s[] %s", header.Name, message.Message)
	} else if header.IP != "" {
		if slices.Contains(cammieBlockedIps, header.IP) {
			c.JSON(http.StatusOK, gin.H{"message": "Message received"})
			return
		}
		newMessage = fmt.Sprintf("<%s> %s", header.IP, message.Message)
	} else {
		newMessage = message.Message
	}

	// Increment messages
	cammieMessages++

	app.Cammie.Update(newMessage)

	c.JSON(http.StatusOK, gin.H{"message": "Message received"})
}
