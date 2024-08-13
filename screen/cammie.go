package screen

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CammieMessage struct {
	Sender  string
	Message string
}

// Component that displays messages received from the website aka cammie chat
type Cammie struct {
	screenApp *ScreenApp
	view      *tview.TextView

	text   string
	buffer string
}

// Initial value, gets adjusted once it's known how much space is available
var maxMessages = 20

// Available colors
var colors = [...]tcell.Color{
	tcell.ColorWhite,
	tcell.ColorViolet,
	tcell.ColorRed,
	tcell.ColorOrange,
	tcell.ColorGreen,
	tcell.ColorAqua,
	tcell.ColorCrimson,
	tcell.ColorFuchsia,
	tcell.ColorGoldenrod,
	tcell.ColorYellow,
	tcell.ColorSalmon,
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
func (cammie *Cammie) Update(message *CammieMessage) {
	colorIndex := hashColor(message.Sender)

	color := colors[colorIndex].String()

	cammie.screenApp.execute(func() {
		fmt.Fprintf(cammie.view, "\n[%s]%s %s[-:-:-:-]", color, message.Sender, message.Message)

		cammie.view.ScrollToEnd()
	})
}

func hashColor(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	hashNumber := h.Sum32()
	return int(hashNumber) % len(colors)
}
