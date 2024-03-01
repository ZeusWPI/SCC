package main

import (
	"scc/api"
	"scc/screen"
)

func main() {
	go api.Start()

	screen.InitApp()
}
