// Package bar provides an animated progress bar
package bar

import (
	"strings"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var lastID int64

func nextID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

// FrameMsg is a message that is sent on every progress frame tick
type FrameMsg struct {
	id int64
}

// StartMsg is a message that starts the progress bar
type StartMsg struct {
	width       int
	widthTarget int
	interval    time.Duration
}

// Model for the progress component
type Model struct {
	id          int64
	width       int
	widthTarget int
	interval    time.Duration
	style       lipgloss.Style
}

// Interface compliance
var _ tea.Model = (*Model)(nil)

func New(style lipgloss.Style) Model {
	return Model{id: nextID(), style: style}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// Start starts a progress bar until it reaches a given width in a given duration
func (m Model) Start(width int, runningTime time.Duration, duration time.Duration) tea.Cmd {
	return func() tea.Msg {
		interval := (duration / 2) / time.Duration(width)

		return StartMsg{
			width:       int(runningTime / interval),
			widthTarget: width * 2,
			interval:    interval,
		}
	}
}

// Update handles the progress frame tick
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FrameMsg:
		if msg.id != m.id {
			return m, nil
		}

		m.width++

		if m.width < m.widthTarget {
			return m, tick(m.id, m.interval)
		}

		return m, nil

	case StartMsg:
		m.id = nextID()
		m.width = msg.width
		m.widthTarget = msg.widthTarget
		m.interval = msg.interval

		return m, tick(m.id, m.interval)
	}

	return m, nil
}

func (m Model) View() string {
	b := strings.Repeat("▄", m.width/2)
	if m.width%2 == 1 {
		b += "▖"
	}
	return m.style.Render(b)
}

func tick(id int64, interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(_ time.Time) tea.Msg {
		return FrameMsg{id: id}
	})
}
