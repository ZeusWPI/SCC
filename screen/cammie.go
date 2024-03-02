package screen

import (
	"scc/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Initial value, gets adjusted once it's known how much space is available
var maxMessages = 20

// Component that displays messages received from the website aka cammie chat
type Cammie struct {
	screenApp *ScreenApp
	view      *tview.TextView

	queue  *utils.Queue[string]
	text   string
	buffer string
}

// Create a new cammie struct
func NewCammie(screenApp *ScreenApp) *Cammie {
	cammie := Cammie{
		screenApp: screenApp,
		view:      tview.NewTextView().SetWrap(true).SetWordWrap(true).SetScrollable(true),

		queue: utils.NewQueue[string](maxMessages),
	}

	cammie.view.SetTitle(" Cammie ")
	cammie.view.SetBorder(true)
	cammie.view.SetTextColor(tcell.ColorOrange)
	cammie.view.SetBorderColor(tcell.ColorOrange)
	cammie.view.SetTitleColor(tcell.ColorOrange)

	return &cammie
}

// One-time setup
func (cammie *Cammie) Run() {
	// Wait for the view to be properly set up
	time.Sleep(5 * time.Second)

	_, _, _, h := cammie.view.GetInnerRect()
	cammie.queue.SetMaxSize(h)
}

// Updates the cammie chat
// Gets called when a new message is received from the website
func (cammie *Cammie) Update(message string) {
	cammie.queue.Enqueue(message)

	cammie.screenApp.execute(func() {
		cammie.view.Clear()

		w := cammie.view.BatchWriter()
		defer w.Close()
		w.Clear()

		for index, message := range cammie.queue.Get() {
			if index != cammie.queue.Size()-1 {
				message += "\n"
			}

			w.Write([]byte(message))
		}

		cammie.view.ScrollToEnd()
	})
}
