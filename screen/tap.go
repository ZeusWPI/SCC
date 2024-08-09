package screen

import (
	"scc/config"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type TapOrder struct {
	OrderID         int    `json:"order_id"`
	OrderCreatedAt  string `json:"order_created_at"`
	ProductName     string `json:"product_name"`
	ProductCategory string `json:"product_category"`
}

type Tap struct {
	ScreenApp *ScreenApp
	view      *tview.Flex
	bar       *tvxwidgets.BarChart
}

var (
	soft = 0
	mate = 0
	beer = 0
)

func NewTap(screenApp *ScreenApp) *Tap {
	tap := Tap{
		ScreenApp: screenApp,
		view:      tview.NewFlex(),
		bar:       tvxwidgets.NewBarChart(),
	}

	tap.view.SetBorder(true).SetTitle(" Tap ")

	tap.bar.AddBar("Soft", 0, tcell.ColorBlue)
	tap.bar.AddBar("Mate", 0, tcell.ColorOrange)
	tap.bar.AddBar("Beer", 0, tcell.ColorRed)
	tap.bar.SetAxesLabelColor(tcell.ColorWhite)

	tap.view.AddItem(tap.bar, 0, 1, true)

	return &tap
}

func (tap *Tap) Run() {
}

func (tap *Tap) Update(order *TapOrder) {
	switch {
	case strings.Contains(order.ProductName, "Mate"):
		mate++
		tap.bar.SetBarValue("Mate", mate)
	case isBeer(order.ProductName):
		beer++
		tap.bar.SetBarValue("Beer", beer)
	default:
		soft++
		tap.bar.SetBarValue("Soft", soft)
	}
}

func isBeer(productName string) bool {
	for _, beer := range config.GetConfig().Tap.Beers {
		if strings.Contains(productName, beer) {
			return true
		}
	}

	return false
}
