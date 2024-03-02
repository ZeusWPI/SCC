package main

import (
	"scc/api"
	"scc/screen"
)

func main() {
	cammieQueue := screen.NewSafeMessageQueue()

	go api.Start(cammieQueue)

	screen.Start(cammieQueue)
}
