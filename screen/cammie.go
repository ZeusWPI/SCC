package screen

import (
	"scc/utils"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Initial value, gets adjusted once it's known how much space is available
var maxMessages = 20

type Cammie struct {
	screenApp *ScreenApp
	view      *tview.TextView

	queue  *utils.Queue[string]
	text   string
	buffer string
}

func NewCammie(screenApp *ScreenApp) *Cammie {
	cammie := Cammie{
		screenApp: screenApp,
		view:      tview.NewTextView().SetWrap(true).SetWordWrap(true).SetText("pls"),

		queue: utils.NewQueue[string](maxMessages),
	}

	cammie.view.SetTitle(" Cammie ")
	cammie.view.SetBorder(true)
	cammie.view.SetTextColor(tcell.ColorOrange)
	cammie.view.SetBorderColor(tcell.ColorOrange)
	cammie.view.SetTitleColor(tcell.ColorOrange)

	return &cammie
}

func (cammie *Cammie) Run() {
	time.Sleep(5 * time.Second)

	_, _, _, h := cammie.view.GetInnerRect()
	cammie.queue.SetMaxSize(h)
}

func (cammie *Cammie) Update(message string) {
	cammie.queue.Enqueue(message)

	cammie.screenApp.execute(func() {
		cammie.screenApp.app.QueueUpdateDraw(func() {
			cammie.view.Clear()

			for _, message := range cammie.queue.Get() {
				cammie.view.Write([]byte(message + "\n"))
			}
		})
	})
}
