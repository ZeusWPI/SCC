package main

import (
	"log"
	"os"
	"scc/api"
	"scc/screen"
)

func main() {
	// Logging
	logFile, err := os.OpenFile("scc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error: Failed to open log file: %s", err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Data holder for the screen
	screenApp := screen.NewScreenApp()

	// Start the API
	go api.Start(screenApp)

	// Start the screen
	screen.Start(screenApp)
}
