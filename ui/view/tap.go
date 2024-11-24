package view

import (
	"context"
	"database/sql"

	"github.com/NimbleMarkets/ntcharts/barchart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
)

// TapModel represents the tap model
type TapModel struct {
	db          *db.DB
	lastOrderID int64
	mate        float64
	soft        float64
	beer        float64
	food        float64
}

// TapMsg represents a tap message
type TapMsg struct {
	lastOrderID int64
	items       []tapItem
}

type tapItem struct {
	category string
	amount   float64
}

var tapCategoryColor = map[string]lipgloss.Color{
	"Mate": lipgloss.Color("208"),
	"Soft": lipgloss.Color("86"),
	"Beer": lipgloss.Color("160"),
	"Food": lipgloss.Color("40"),
}

// NewTapModel creates a new tap model
func NewTapModel(db *db.DB) View {
	return &TapModel{db: db, lastOrderID: -1}
}

// Init initializes the tap model
func (t *TapModel) Init() tea.Cmd {
	return nil
}

// Update updates the tap model
func (t *TapModel) Update(msg tea.Msg) (View, tea.Cmd) {
	switch msg := msg.(type) {
	case TapMsg:
		t.lastOrderID = msg.lastOrderID

		for _, msg := range msg.items {
			switch msg.category {
			case "Mate":
				t.mate += msg.amount
			case "Soft":
				t.soft += msg.amount
			case "Beer":
				t.beer += msg.amount
			case "Food":
				t.food += msg.amount
			}
		}

		return t, nil
	}

	return t, nil
}

// View returns the tap view
func (t *TapModel) View() string {
	chart := barchart.New(20, 20)

	barMate := barchart.BarData{
		Label: "Mate",
		Values: []barchart.BarValue{{
			Name:  "Mate",
			Value: t.mate,
			Style: lipgloss.NewStyle().Foreground(tapCategoryColor["Mate"]),
		}},
	}
	barSoft := barchart.BarData{
		Label: "Soft",
		Values: []barchart.BarValue{{
			Name:  "Soft",
			Value: t.soft,
			Style: lipgloss.NewStyle().Foreground(tapCategoryColor["Soft"]),
		}},
	}
	barBeer := barchart.BarData{
		Label: "Beer",
		Values: []barchart.BarValue{{
			Name:  "Beer",
			Value: t.beer,
			Style: lipgloss.NewStyle().Foreground(tapCategoryColor["Beer"]),
		}},
	}
	barFood := barchart.BarData{
		Label: "Food",
		Values: []barchart.BarValue{{
			Name:  "Food",
			Value: t.food,
			Style: lipgloss.NewStyle().Foreground(tapCategoryColor["Food"]),
		}},
	}

	chart.PushAll([]barchart.BarData{barMate, barSoft, barBeer, barFood})
	chart.Draw()

	return chart.View()
}

// GetUpdateDatas returns all the update functions for the tap model
func (t *TapModel) GetUpdateDatas() []UpdateData {
	return []UpdateData{
		{
			Name:     "tap orders",
			View:     t,
			Update:   updateOrders,
			Interval: config.GetDefaultInt("tui.tap.interval_s", 60),
		},
	}
}

func updateOrders(db *db.DB, view View) (tea.Msg, error) {
	t := view.(*TapModel)
	lastOrderID := t.lastOrderID

	order, err := db.Queries.GetLastOrderByOrderID(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return TapMsg{lastOrderID: lastOrderID, items: []tapItem{}}, err
	}

	if order.OrderID <= lastOrderID {
		return TapMsg{lastOrderID: lastOrderID, items: []tapItem{}}, nil
	}

	orders, err := db.Queries.GetOrderCountByCategorySinceOrderID(context.Background(), lastOrderID)
	if err != nil {
		return TapMsg{lastOrderID: lastOrderID, items: []tapItem{}}, err
	}

	mate, soft, beer, food := 0.0, 0.0, 0.0, 0.0
	for _, order := range orders {
		switch order.Category {
		case "Mate":
			mate += float64(order.Count)
		case "Soft":
			soft += float64(order.Count)
		case "Beer":
			beer += float64(order.Count)
		case "Food":
			food += float64(order.Count)
		}
	}

	messages := make([]tapItem, 0, 4)
	if mate > 0 {
		messages = append(messages, tapItem{"Mate", mate})
	}
	if soft > 0 {
		messages = append(messages, tapItem{"Soft", soft})
	}
	if beer > 0 {
		messages = append(messages, tapItem{"Beer", beer})
	}
	if food > 0 {
		messages = append(messages, tapItem{"Food", food})
	}

	return TapMsg{lastOrderID: order.OrderID, items: messages}, err
}
