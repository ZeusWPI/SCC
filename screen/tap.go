package screen

import (
	"scc/config"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type TapOrder struct {
	OrderID         int       `json:"order_id"`
	OrderCreatedAt  time.Time `json:"order_created_at"`
	ProductName     string    `json:"product_name"`
	ProductCategory string    `json:"product_category"`
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
	food = 0
)

func NewTap(screenApp *ScreenApp) *Tap {
	tap := Tap{
		ScreenApp: screenApp,
		view:      tview.NewFlex(),
		bar:       tvxwidgets.NewBarChart(),
	}

	tap.view.SetBorder(true).SetTitle(" Tap ")

	tap.bar.AddBar("Soft", 0, tcell.ColorAqua)
	tap.bar.AddBar("Mate", 0, tcell.ColorOrange)
	tap.bar.AddBar("Beer", 0, tcell.ColorRed)
	tap.bar.AddBar("Food", 0, tcell.ColorGreen)
	tap.bar.SetAxesColor(tcell.ColorWhite)
	tap.bar.SetAxesLabelColor(tcell.ColorWhite)

	tap.view.AddItem(tap.bar, 0, 1, false)

	return &tap
}

func (tap *Tap) Run() {
}

func (tap *Tap) Update(order *TapOrder) {
	var label string
	var value *int

	switch {
	case order.ProductCategory == "food":
		label, value = "Food", &food
	case order.ProductCategory != "beverages":
		return
	case strings.Contains(order.ProductName, "Mate"):
		label, value = "Mate", &mate
	case isBeer(order.ProductName):
		label, value = "Beer", &beer
	default:
		label, value = "Soft", &soft
	}

	*value++
	tap.ScreenApp.execute(func() {
		tap.bar.SetBarValue(label, *value)
	})
}

func isBeer(productName string) bool {
	for _, beer := range config.GetConfig().Tap.Beers {
		if strings.Contains(productName, beer) {
			return true
		}
	}

	return false
}
