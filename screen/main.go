package screen

import (
	"sync"

	"github.com/rivo/tview"
)

// Main struct for the screen application
type ScreenApp struct {
	mu  sync.Mutex
	app *tview.Application

	Spotify *Spotify
	Cammie  *Cammie
	Graph1  *Graph1
	Graph2  *Graph2
}

// Execute a function with a lock
func (screenApp *ScreenApp) execute(f func()) {
	screenApp.mu.Lock()
	defer screenApp.mu.Unlock()
	f()
}

// Create a new screen application
func NewScreenApp() *ScreenApp {
	screen := ScreenApp{
		app: tview.NewApplication(),
	}

	screen.Spotify = NewSpotify(&screen)
	screen.Cammie = NewCammie(&screen)
	screen.Graph1 = NewGraph1(&screen)
	screen.Graph2 = NewGraph2(&screen)

	// Build the screen layout
	screen.app.SetRoot(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(screen.Spotify.view, 3, 2, false).
		AddItem(tview.NewFlex().
			AddItem(screen.Cammie.view, 0, 5, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(screen.Graph1.view, 0, 1, false).
				AddItem(screen.Graph2.view, 0, 1, false), 0, 4, false), 0, 13, false), true).
		EnableMouse(true)

	return &screen
}

// Start the screen application
func Start(screen *ScreenApp) {

	// Start each screen component
	go screen.Spotify.Run()
	go screen.Cammie.Run()
	go screen.Graph1.Run()
	go screen.Graph2.Run()

	// Start the screen application
	if err := screen.app.Run(); err != nil {
		panic(err)
	}
}
