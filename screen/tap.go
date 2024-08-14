package screen

import (
	"scc/config"
	"sort"
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
type tapItem struct {
	value int
	color tcell.Color
}

const (
	soft = "Soft"
	mate = "Mate"
	beer = "Beer"
	food = "Food"
)

var tapItems = map[string]tapItem{
	soft: {0, tcell.ColorAqua},
	mate: {0, tcell.ColorOrange},
	beer: {0, tcell.ColorRed},
	food: {0, tcell.ColorGreen},
}

func NewTap(screenApp *ScreenApp) *Tap {
	tap := Tap{
		ScreenApp: screenApp,
		view:      tview.NewFlex(),
		bar:       tvxwidgets.NewBarChart(),
	}

	tap.view.SetBorder(true).SetTitle(" Tap ")

	tap.bar.SetAxesColor(tcell.ColorWhite)
	tap.bar.SetAxesLabelColor(tcell.ColorWhite)

	tap.view.AddItem(tap.bar, 0, 1, false)

	return &tap
}

func (tap *Tap) Run() {
}

func (tap *Tap) Update(order *TapOrder) {
	var key string
	switch {
	case order.ProductCategory == "food":
		key = food
	case order.ProductCategory != "beverages":
		return
	case strings.Contains(order.ProductName, "Mate"):
		key = mate
	case isBeer(order.ProductName):
		key = beer
	default:
		key = soft
	}

	entry := tapItems[key]
	entry.value++
	tapItems[key] = entry

	// item.amount++
	tap.ScreenApp.execute(func() {
		// Remove labels
		for label := range tapItems {
			tap.bar.RemoveBar(label)
		}

		// Create slice of keys
		keys := make([]string, 0, len(tapItems))
		for key := range tapItems {
			keys = append(keys, key)
		}

		// Sort slice
		sort.Slice(keys, func(i, j int) bool {
			return tapItems[keys[i]].value > tapItems[keys[j]].value
		})

		// Add labels back
		for _, key := range keys {
			tap.bar.AddBar(key, tapItems[key].value, tapItems[key].color)
		}

		// Required so that the bars change relative height
		tap.bar.SetMaxValue(tapItems[keys[0]].value)
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
