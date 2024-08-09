package screen

import "github.com/rivo/tview"

type Tap struct {
	ScreenApp *ScreenApp
	view      *tview.Box
}

func NewTap(screenApp *ScreenApp) *Tap {
	tap := Tap{
		ScreenApp: screenApp,
		view:      tview.NewBox().SetBorder(true).SetTitle(" Tap "),
	}

	return &tap
}

func (tap *Tap) Run() {
}

func (tap *Tap) Update(text string) {
}
