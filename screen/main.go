package screen

import (
	"sync"

	"github.com/rivo/tview"
)

type ScreenApp struct {
	mu  sync.Mutex
	app *tview.Application

	Spotify *Spotify
	Cammie  *Cammie
	Graph1  *Graph1
	Graph2  *Graph2
}

type s_cammie struct {
	cammie *tview.TextView
}

func NewScreenApp() *ScreenApp {
	screen := ScreenApp{
		app: tview.NewApplication(),
	}

	screen.Spotify = NewSpotify(&screen)
	screen.Cammie = NewCammie(&screen)
	screen.Graph1 = NewGraph1(&screen)
	screen.Graph2 = NewGraph2(&screen)

	screen.app.SetRoot(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(screen.Spotify.view, 3, 2, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true).SetBorder(true).SetTitle("Cammie"), 0, 5, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Graph 1"), 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Graph 2"), 0, 1, false), 0, 4, false), 0, 13, false), true).
		EnableMouse(true)

	return &screen
}

func Start(screen *ScreenApp) {

	go screen.Spotify.Run()
	go screen.Cammie.Run()
	go screen.Graph1.Run()
	go screen.Graph2.Run()

	if err := screen.app.Run(); err != nil {
		panic(err)
	}
}
