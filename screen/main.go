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
	// spotify *tview.TextView
	// cammie  *s_cammie
	// graph1  *tview.Box
	// graph2  *tview.Box
}

type s_cammie struct {
	cammie *tview.TextView
}

func NewScreenApp() *ScreenApp {
	screen := ScreenApp{
		app: tview.NewApplication(),
		// cammie: &s_cammie{
		// 	cammie: tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true),
		// 	queue:  queue,
		// },
		// graph1: tview.NewBox().SetBorder(true).SetTitle("Graph 1"),
		// graph2: tview.NewBox().SetBorder(true).SetTitle("Graph 2"),
	}

	screen.Spotify = NewSpotify(&screen)
	screen.Cammie = NewCammie(&screen)

	screen.app.SetRoot(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(screen.Spotify.view, 3, 2, false).
		AddItem(tview.NewFlex().
			AddItem(screen.Cammie.view, 0, 5, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(screen.graph1, 0, 1, false).
				AddItem(screen.graph2, 0, 1, false), 0, 4, false), 0, 13, false), true).
		EnableMouse(true)

	return &screen
}

func Start(screen *ScreenApp) {

	go screen.Spotify.Run()
	go screen.Cammie.Run()

	if err := screen.app.Run(); err != nil {
		panic(err)
	}
}
