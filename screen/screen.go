package screen

import "github.com/rivo/tview"

func InitApp() {
	app := tview.NewApplication()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Spotify"), 0, 2, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Cammie"), 0, 5, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Graph 1"), 0, 1, false).
				AddItem(tview.NewBox().SetBorder(true).SetTitle("Graph 2"), 0, 1, false), 0, 4, false), 0, 13, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}