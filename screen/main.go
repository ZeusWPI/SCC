package screen

import (
	"sync"

	"github.com/rivo/tview"
)

type screenApp struct {
	mu      sync.Mutex
	app     *tview.Application
	spotify *tview.Box
	cammie  *s_cammie
	graph1  *tview.Box
	graph2  *tview.Box
}

type s_cammie struct {
	cammie *tview.TextView
	queue  *SafeMessageQueue
}

var screen screenApp

func Start(queue *SafeMessageQueue) {
	screen = screenApp{
		app:     tview.NewApplication(),
		spotify: tview.NewBox().SetBorder(true).SetTitle("Spotify"),
		cammie: &s_cammie{
			cammie: tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true),
			queue:  queue,
		},
		graph1: tview.NewBox().SetBorder(true).SetTitle("Graph 1"),
		graph2: tview.NewBox().SetBorder(true).SetTitle("Graph 2"),
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(screen.spotify, 0, 2, false).
		AddItem(tview.NewFlex().
			AddItem(screen.cammie.cammie.SetBorder(true).SetTitle("Cammie"), 0, 5, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(screen.graph1, 0, 1, false).
				AddItem(screen.graph2, 0, 1, false), 0, 4, false), 0, 13, false)

	go updateCammie()

	if err := screen.app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
