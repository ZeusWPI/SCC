package main

import (
	"scc/api"
	"scc/screen"
)

func main() {
	cammieQueue := screen.NewSafeMessageQueue()

	screenApp := screen.NewScreenApp()

	go api.Start(screenApp, cammieQueue)

	screen.Start(screenApp, cammieQueue)
}
