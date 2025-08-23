// Package event provides the functions to draw all the upcoming zeus events on a TUI
package event

import (
	"image"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

type Model struct {
	events []event

	width  int
	height int

	url string // Url of the api
}

// Interface compliance
var _ view.View = (*Model)(nil)

// Msg represents the message to update the event view
type Msg struct {
	events []event
}

// event is the internal representation of a zeus event
type event struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Start    time.Time `json:"start_time"`
	poster   image.Image
}

func NewModel() view.View {
	return &Model{
		events: nil,
		width:  0,
		height: 0,
		url:    config.GetDefaultString("tui.view.event.url", "https://events.zeus.gent/api/v1"),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Name() string {
	return "Events"
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
		m.events = msg.events
	}

	return m, nil
}

func (m *Model) View() string {
	if idx := slices.IndexFunc(m.events, func(e event) bool { return utils.SameDay(e.Start, time.Now()) }); idx != -1 {
		return m.viewToday()
	}

	return m.viewOverview()
}

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
