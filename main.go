package main

import (
	"scc/api"
	"scc/screen"
)

func main() {
	// Data holder for the screen
	screenApp := screen.NewScreenApp()

	// Start the API
	go api.Start(screenApp)

	// Start the screen
	screen.Start(screenApp)
}
