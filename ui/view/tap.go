// Package view contains all the different views for the tui
package view

import (
	"context"
	"time"

	"github.com/NimbleMarkets/ntcharts/barchart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"go.uber.org/zap"
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

// TapMessage represents a tap message
type TapMessage struct {
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
func NewTapModel(db *db.DB) *TapModel {
	return &TapModel{db: db, lastOrderID: -1}
}

// Init initializes the tap model
func (t *TapModel) Init() tea.Cmd {
	return updateOrders(t.db, t.lastOrderID)
}

// Update updates the tap model
func (t *TapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TapMessage:
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

		return t, updateOrders(t.db, t.lastOrderID)
	}

	return t, nil
}

// View returns the tap view
func (t *TapModel) View() string {
	chart := barchart.New(20, 20)

	barMate := barchart.BarData{
		Label: "Mate",
		Values: []barchart.BarValue{{
			Name:  "Item1",
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

func updateOrders(db *db.DB, lastOrderID int64) tea.Cmd {
	return tea.Tick(60*time.Second, func(_ time.Time) tea.Msg {
		order, err := db.Queries.GetLastOrderByOrderID(context.Background())
		if err != nil {
			zap.S().Error("DB: Failed to get last order", err)
			return TapMessage{lastOrderID: lastOrderID, items: nil}
		}

		if order.OrderID <= lastOrderID {
			return TapMessage{lastOrderID: lastOrderID, items: nil}
		}

		orders, err := db.Queries.GetOrderCountByCategorySinceOrderID(context.Background(), lastOrderID)
		if err != nil {
			zap.S().Error("DB: Failed to get tap orders", err)
			return TapMessage{lastOrderID: lastOrderID, items: nil}
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

		return TapMessage{lastOrderID: order.OrderID, items: messages}
	})
}
