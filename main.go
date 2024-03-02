package main

import (
	"scc/api"
	"scc/screen"
)

func main() {
	screenApp := screen.NewScreenApp()

	go api.Start(screenApp)

	screen.Start(screenApp)
}
