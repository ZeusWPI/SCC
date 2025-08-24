// Package zess provides the functions to draw an overview of the zess scans on a TUI
package zess

import (
	"math/rand/v2"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/view"
)

type Model struct {
	repoSeason repository.Season
	repoScan   repository.Scan

	weeks  []week
	width  int
	height int
}

// Interface compliance
var _ view.View = (*Model)(nil)

type Msg struct {
	weeks []week
}

type week struct {
	start time.Time // Start of the week

	scans int
}

func NewModel(repo repository.Repository) view.View {
	m := &Model{
		repoSeason: *repo.NewSeason(),
		repoScan:   *repo.NewScan(),
		weeks:      nil,
		width:      0,
		height:     0,
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Name() string {
	return "Zess"
}

func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		entry, ok := msg.Sizes[m.Name()]
		if ok {
			// Update all dependent styles
			m.width = entry.Width
			m.height = entry.Height

			m.updateStyles()
		}

		return m, nil

	case Msg:
		m.weeks = msg.weeks
	}

	return m, nil
}

// View returns the view for the zess model
func (m *Model) View() string {
	chart := m.viewChart()
	stats := m.viewStats()

	// Join them together
	view := lipgloss.JoinHorizontal(lipgloss.Top, chart, stats)
	return view
}

// GetUpdateDatas returns all the update functions for the zess model
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "zess weeks",
			View:     m,
			Update:   updateWeeks,
			Interval: config.GetDefaultInt("tui.view.zess.interval_s", 60),
		},
	}
}

func randomColor() string {
	return colors[rand.IntN(len(colors))]
}
