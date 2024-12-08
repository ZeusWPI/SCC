// Package tap provides the functions to draw an overview of the recent tap orders on a TUI
package tap

import (
	"context"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jackc/pgx/v5"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/view"
)

type category string

const (
	mate category = "Mate"
	soft category = "Soft"
	beer category = "Beer"
	food category = "Food"
)

var categoryToStyle = map[category]lipgloss.Style{
	mate: sMate,
	soft: sSoft,
	beer: sBeer,
	food: sFood,
}

// Model represents the tap model
type Model struct {
	db          *db.DB
	lastOrderID int32
	items       []tapItem
}

// Msg represents a tap message
type Msg struct {
	lastOrderID int32
	items       []tapItem
}

type tapItem struct {
	category category
	amount   int
	last     time.Time
}

// NewModel creates a new tap model
func NewModel(db *db.DB) view.View {
	return &Model{db: db, lastOrderID: -1}
}

// Init initializes the tap model
func (m *Model) Init() tea.Cmd {
	return nil
}

// Name returns the name of the view
func (m *Model) Name() string {
	return "Tap"
}

// Update updates the tap model
func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		m.lastOrderID = msg.lastOrderID

		for _, msgItem := range msg.items {
			found := false
			for i, item := range m.items {
				if item.category == msgItem.category {
					m.items[i].amount += msgItem.amount
					m.items[i].last = msgItem.last
					found = true
					break
				}
			}

			if !found {
				m.items = append(m.items, msgItem)
			}
		}

		// Sort to display bars in order
		slices.SortFunc(m.items, func(i, j tapItem) int {
			return j.amount - i.amount
		})

		return m, nil
	}

	return m, nil
}

// View returns the tap view
func (m *Model) View() string {
	chart := m.viewChart()
	stats := m.viewStats()

	// Give them same height
	stats = sStats.Height(lipgloss.Height(chart)).Render(stats)

	// Join them together
	view := lipgloss.JoinHorizontal(lipgloss.Top, chart, stats)
	return view
}

// GetUpdateDatas returns all the update functions for the tap model
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "tap orders",
			View:     m,
			Update:   updateOrders,
			Interval: config.GetDefaultInt("tui.tap.interval_s", 60),
		},
	}
}

func updateOrders(view view.View) (tea.Msg, error) {
	m := view.(*Model)
	lastOrderID := m.lastOrderID

	order, err := m.db.Queries.GetLastOrderByOrderID(context.Background())
	if err != nil {
		if err == pgx.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	if order.OrderID <= lastOrderID {
		return nil, nil
	}

	orders, err := m.db.Queries.GetOrderCountByCategorySinceOrderID(context.Background(), lastOrderID)
	if err != nil {
		return nil, err
	}

	counts := make(map[category]tapItem)

	for _, order := range orders {
		if entry, ok := counts[category(order.Category)]; ok {
			entry.amount += int(order.Count)
			counts[category(order.Category)] = entry
			continue
		}

		counts[category(order.Category)] = tapItem{
			category: category(order.Category),
			amount:   int(order.Count),
			last:     time.Unix(int64(order.LatestOrderCreatedAt), 0),
		}
	}

	items := make([]tapItem, 0, len(counts))

	for _, v := range counts {
		items = append(items, v)
	}

	return Msg{lastOrderID: order.OrderID, items: items}, nil
}
