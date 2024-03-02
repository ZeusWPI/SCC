package screen

import "github.com/rivo/tview"

type Graph2 struct {
	ScreenApp *ScreenApp
	view      *tview.Box
}

func NewGraph2(screenApp *ScreenApp) *Graph2 {
	graph2 := Graph2{
		ScreenApp: screenApp,
		view:      tview.NewBox().SetBorder(true).SetTitle("Graph 1"),
	}

	return &graph2
}

func (graph1 *Graph2) Run() {
}

func (graph1 *Graph2) Update(text string) {
}
