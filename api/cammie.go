package api

import (
	"fmt"
	"net/http"
	"scc/screen"
	"slices"

	gin "github.com/gin-gonic/gin"
)

// messageCammie struct
type messageCammie struct {
	Message string `form:"message" json:"message" xml:"message" binding:"required"`
}

// headerCammie struct
type headerCammie struct {
	Name string `header:"X-Username"`
	IP   string `header:"X-Real-IP"`
}

var messages uint64 = 0
var blockedNames = []string{"Paul-Henri Spaak"} // Blocekd names
var blockedIps = []string{}                     // Blocked IPs
var maxMessageLength = 200                      // Maximum message length

func getMessage(app *screen.ScreenApp, c *gin.Context) {
	c.JSON(200, gin.H{"messages": messages})
}

func postMessage(app *screen.ScreenApp, c *gin.Context) {
	// Get structs
	header := &headerCammie{}
	message := &messageCammie{}

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
		newMessage = fmt.Sprintf("[%s[] %s", header.Name, message.Message)
	} else if header.IP != "" {
		if slices.Contains(blockedIps, header.IP) {
			c.JSON(http.StatusOK, gin.H{"message": "Message received"})
			return
		}
		newMessage = fmt.Sprintf("<%s> %s", header.IP, message.Message)
	} else {
		newMessage = message.Message
	}

	// Increment messages
	messages++

	app.Cammie.Update(newMessage)

	c.JSON(http.StatusOK, gin.H{"message": "Message received"})
}
