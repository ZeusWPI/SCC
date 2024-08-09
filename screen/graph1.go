package screen

import "github.com/rivo/tview"

type Graph1 struct {
	ScreenApp *ScreenApp
	view      *tview.Box
}

func NewGraph1(screenApp *ScreenApp) *Graph1 {
	graph1 := Graph1{
		ScreenApp: screenApp,
		view:      tview.NewBox().SetBorder(true).SetTitle(" Graph 1 "),
	}

	return &graph1
}

func (graph1 *Graph1) Run() {
}

func (graph1 *Graph1) Update(text string) {
}
