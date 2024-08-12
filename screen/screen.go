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
	Tap     *Tap
	Zess    *Zess
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
	screen.Tap = NewTap(&screen)
	screen.Zess = NewZess(&screen)

	// Build the screen layout
	screen.app.SetRoot(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(screen.Spotify.view, 3, 2, false).
		AddItem(tview.NewFlex().
			AddItem(screen.Cammie.view, 0, 5, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(screen.Tap.view, 0, 1, false).
				AddItem(screen.Zess.view, 0, 1, false), 0, 4, false), 0, 13, false), true).
		EnableMouse(false)

	return &screen
}

// Start the screen application
func Start(screen *ScreenApp) {

	// Start each screen component
	go screen.Spotify.Run()
	go screen.Cammie.Run()
	go screen.Tap.Run()
	go screen.Zess.Run()

	// Start the screen application
	if err := screen.app.Run(); err != nil {
		panic(err)
	}
}
