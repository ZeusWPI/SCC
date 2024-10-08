package screen

import (
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Spotify struct {
	screenApp *ScreenApp
	view      *tview.TextView

	mu     sync.Mutex
	text   string
	buffer string
}

func NewSpotify(screenApp *ScreenApp) *Spotify {
	spotify := Spotify{
		screenApp: screenApp,
		view:      tview.NewTextView(),

		text:   "VERY COOL SONG - Le Artist",
		buffer: "",
	}

	spotify.view.SetTitle(" Spotify ")
	spotify.view.SetBorder(true)
	spotify.view.SetTextColor(tcell.ColorGreen)
	spotify.view.SetBorderColor(tcell.ColorGreen)
	spotify.view.SetTitleColor(tcell.ColorGreen)

	return &spotify
}

func (spotify *Spotify) Run() {
	time.Sleep(1 * time.Second)

	for {
		_, _, w, _ := spotify.view.GetInnerRect()

		if w != 0 {

			spotify.screenApp.execute(func() {
				if len(spotify.buffer) != w {
					if len(spotify.text) > w {
						spotify.text = spotify.text[0 : w-4]
						spotify.text += "..."
					}
					spotify.buffer = spotify.text + strings.Repeat(" ", w-len(spotify.text))
				}

				spotify.buffer = spotify.buffer[1:] + string(spotify.buffer[0])
			})

			spotify.screenApp.app.QueueUpdateDraw(func() {
				spotify.view.SetText(spotify.buffer)
			})
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (spotify *Spotify) Update(text string) {
	spotify.screenApp.execute(func() {
		spotify.text = text
		spotify.buffer = ""
	})
}
