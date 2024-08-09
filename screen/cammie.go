package screen

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Initial value, gets adjusted once it's known how much space is available
var maxMessages = 20

// Available colors
var colors = [...]tcell.Color{
	tcell.ColorViolet,
	tcell.ColorRed,
	tcell.ColorIndigo,
	tcell.ColorYellow,
	tcell.ColorBlue,
	tcell.ColorGreen,
	tcell.ColorOrange,
	tcell.ColorLime,
	tcell.ColorAqua,
	tcell.ColorDarkSalmon,
	tcell.ColorLightBlue,
	tcell.ColorNavajoWhite,
}
var lastColorIndex = 0

// Component that displays messages received from the website aka cammie chat
type Cammie struct {
	screenApp *ScreenApp
	view      *tview.TextView

	text   string
	buffer string
}

// Create a new cammie struct
func NewCammie(screenApp *ScreenApp) *Cammie {
	cammie := Cammie{
		screenApp: screenApp,
		view: tview.NewTextView().
			SetWordWrap(true).
			SetScrollable(true).
			SetDynamicColors(true),
	}

	cammie.view.SetTitle(" Cammie ").
		SetBorder(true).
		SetBorderColor(tcell.ColorOrange).
		SetTitleColor(tcell.ColorOrange)

	return &cammie
}

// Run one-time setup
func (cammie *Cammie) Run() {
	// Wait for the view to be properly set up
	time.Sleep(1 * time.Second)
}

// Updates the cammie chat
// Gets called when a new message is received from the website
func (cammie *Cammie) Update(message string) {
	color := colors[lastColorIndex].String()
	lastColorIndex = (lastColorIndex + 1) % len(colors)

	cammie.screenApp.execute(func() {
		fmt.Fprintf(cammie.view, "\n[%s]%s", color, message)

		cammie.view.ScrollToEnd()
	})
}
