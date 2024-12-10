// Package stopwatch provides a simple stopwatch component
package stopwatch

import (
	"fmt"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Slightly adjusted version of https://github.com/charmbracelet/bubbles/blob/master/stopwatch/stopwatch.go

var lastID int64

func nextID() int64 {
	return atomic.AddInt64(&lastID, 1)
}

// TickMsg is a message that is sent on every stopwatch tick
type TickMsg struct {
	id int64
}

// StartStopMsg is a message that controls if the stopwatch is running or not
type StartStopMsg struct {
	running       bool
	startDuration time.Duration
}

// ResetMsg is a message that resets the stopwatch
type ResetMsg struct {
}

// Model for the stopwatch component
type Model struct {
	id       int64
	duration time.Duration
	running  bool
}

// New creates a new stopwatch with a given interval
func New() Model {
	return Model{
		id:       nextID(),
		duration: 0,
		running:  false,
	}
}

// Init initializes the stopwatch component
func (m Model) Init() tea.Cmd {
	return nil
}

// Start starts the stopwatch
func (m Model) Start(startDuration time.Duration) tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{running: true, startDuration: startDuration}
	}
}

// Stop stops the stopwatch
func (m Model) Stop() tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{running: false}
	}
}

// Reset resets the stopwatch
func (m Model) Reset() tea.Cmd {
	return func() tea.Msg {
		return ResetMsg{}
	}
}

// Update handles the stopwatch tick
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if msg.id != m.id || !m.running {
			return m, nil
		}

		m.duration += time.Second
		return m, tick(m.id)

	case ResetMsg:
		m.duration = 0
		m.running = false

	case StartStopMsg:
		if msg.running {
			// Start
			if m.running {
				// Already running
				return m, nil
			}

			m.id = nextID()
			m.duration = msg.startDuration
			m.running = true

			return m, tick(m.id)
		}

		// Stop
		m.running = false
		return m, nil
	}
	return m, nil
}

// View of the stopwatch component
func (m Model) View() string {
	duration := m.duration.Round(time.Second)

	min := int(duration / time.Minute)
	sec := int((duration % time.Minute) / time.Second)

	return fmt.Sprintf("%02d:%02d", min, sec)
}

func tick(id int64) tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return TickMsg{id: id}
	})
}
