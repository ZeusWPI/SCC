// Package event provides the functions to draw all the upcoming zeus events on a TUI
package event

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/util"
	"github.com/zeusWPI/scc/tui/view"
)

var (
	passedAmount   = 3
	upcomingAmount = 7
)

// Model represents the model for the event view
type Model struct {
	db       *db.DB
	passed   []dto.Event
	upcoming []dto.Event
	today    *dto.Event

	width  int
	height int
}

// Msg represents the message to update the event view
type Msg struct {
	upcoming []dto.Event
	passed   []dto.Event
	today    *dto.Event
}

// NewModel creates a new event view
func NewModel(db *db.DB) view.View {
	return &Model{db: db}
}

// Init initializes the event model view
func (m *Model) Init() tea.Cmd {
	return nil
}

// Name returns the name of the view
func (m *Model) Name() string {
	return "Events"
}

// Update updates the event model view
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
		m.passed = msg.passed
		m.upcoming = msg.upcoming
		m.today = msg.today
	}

	return m, nil
}

// View returns the view for the event model
func (m *Model) View() string {
	if m.today != nil {
		return m.viewToday()
	}

	return m.viewOverview()
}

// GetUpdateDatas returns all the update function for the event model
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "event update",
			View:     m,
			Update:   updateEvents,
			Interval: config.GetDefaultInt("tui.view.event.interval_s", 3600),
		},
	}
}

func updateEvents(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	eventsDB, err := m.db.Queries.GetEventsCurrentAcademicYear(context.Background())
	if err != nil {
		return nil, err
	}

	events := util.SliceMap(eventsDB, dto.EventDTO)

	passed := make([]dto.Event, 0)
	upcoming := make([]dto.Event, 0)
	var today *dto.Event

	now := time.Now()
	for _, event := range events {
		if event.Date.Before(now) {
			passed = append(passed, *event)
		} else {
			upcoming = append(upcoming, *event)
		}

		if event.Date.Year() == now.Year() && event.Date.YearDay() == now.YearDay() {
			today = event
		}
	}

	// Truncate passed and upcoming slices
	if len(passed) > passedAmount {
		passed = passed[len(passed)-passedAmount:]
	}

	if len(upcoming) > upcomingAmount {
		upcoming = upcoming[:upcomingAmount]
	}

	return Msg{passed: passed, upcoming: upcoming, today: today}, nil
}
