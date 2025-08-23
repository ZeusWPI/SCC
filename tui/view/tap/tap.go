// Package tap provides the functions to draw an overview of the recent tap orders on a TUI
package tap

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/view"
)

var categoryToStyle = map[model.TapCategory]lipgloss.Style{
	model.Mate: sMate,
	model.Soft: sSoft,
	model.Beer: sBeer,
	model.Food: sFood,
}

type Model struct {
	repo        repository.Tap
	lastOrderID int
	items       []model.TapCount

	width  int
	height int
}

// Interface compliance
var _ view.View = (*Model)(nil)

// Msg represents a tap message
type Msg struct {
	lastOrderID int
	items       []model.TapCount
}

func NewModel(repo repository.Repository) view.View {
	return &Model{
		repo:        *repo.NewTap(),
		lastOrderID: -1,
		items:       nil,
		width:       0,
		height:      0,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Name() string {
	return "Tap"
}

func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		if entry, ok := msg.Sizes[m.Name()]; ok {
			// Update all dependent styles
			m.width = entry.Width
			m.height = entry.Height

			m.updateStyles()
		}

		return m, nil

	case Msg:
		m.lastOrderID = msg.lastOrderID
		m.items = msg.items
	}

	return m, nil
}

// View returns the tap view
func (m *Model) View() string {
	chart := m.viewChart()
	stats := m.viewStats()

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
			Interval: config.GetDefaultInt("tui.view.tap.interval_s", 60),
		},
	}
}
