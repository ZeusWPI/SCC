package screen

import "github.com/rivo/tview"

type Cammie struct {
	ScreenApp *ScreenApp
	view      *tview.TextView
}

func NewCammie(screenApp *ScreenApp) *Cammie {
	cammie := Cammie{
		ScreenApp: screenApp,
		view:      tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true),
	}

	cammie.view.SetTitle(" Cammie ")
	cammie.view.SetBorder(true)
	cammie.view.SetTextColor(tview.Styles.PrimaryTextColor)
	cammie.view.SetBorderColor(tview.Styles.BorderColor)

	return &cammie
}

func (cammie *Cammie) Run() {
}

func (cammie *Cammie) Update(text string) {
	cammie.view.SetText(text)
}
