package screen

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Spotify struct {
	screenApp *ScreenApp
	app       *tview.Application
	view      *tview.TextView
}

func NewSpotify(screenApp *ScreenApp) *Spotify {
	spotify := Spotify{
		screenApp: screenApp,
		view:      tview.NewTextView(),
	}

	spotify.view.SetTitle(" Spotify ")
	spotify.view.SetBorder(true)
	spotify.view.SetTextColor(tcell.ColorLimeGreen)
	spotify.view.SetBorderColor(tcell.ColorLimeGreen)

	return &spotify
}

func (spotify *Spotify) Run() {
	i := 0
	for {
		time.Sleep(1 * time.Second)
		spotify.app.QueueUpdateDraw(func() {
			spotify.view.SetText(fmt.Sprintf("%d", i))
		})
		i++
	}
}

func (spotify *Spotify) Update(text string) {

}
