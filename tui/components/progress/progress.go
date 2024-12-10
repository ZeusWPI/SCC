// Package progress provides an animated progress bar
package progress

import (
	"strings"
	"sync/atomic"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
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
	id           int64
	width        int
	widthTarget  int
	interval     time.Duration
	styleFainted lipgloss.Style
	styleGlow    lipgloss.Style
}

// New creates a new progress
func New(styleFainted, styleGlow lipgloss.Style) Model {
	zap.S().Info(styleFainted)
	return Model{id: nextID(), styleFainted: styleFainted, styleGlow: styleGlow}
}

// Init initializes the progress component
func (m Model) Init() tea.Cmd {
	zap.S().Info(m.styleFainted)
	return nil
}

// Start starts a progress bar until it reaches a given width in a given duration
func (m Model) Start(width int, runningTime time.Duration, duration time.Duration) tea.Cmd {
	zap.S().Info(m.styleFainted)
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
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	zap.S().Info(m.styleFainted)
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

// View of the progress bar component
func (m Model) View() string {
	zap.S().Info(m.styleFainted)
	glowCount := min(20, m.width)
	// Make sure if m.width is uneven that the half block string is in the glow part
	if m.width%2 == 1 && glowCount%2 == 0 {
		glowCount--
	}
	faintedCount := m.width - glowCount

	// Construct fainted
	fainted := strings.Repeat("▄", faintedCount/2)
	fainted = m.styleFainted.Render(fainted)

	// Construct glow
	glow := strings.Repeat("▄", glowCount/2)
	if glowCount%2 == 1 {
		glow += "▖"
	}
	glow = m.styleGlow.Render(glow)

	return lipgloss.JoinHorizontal(lipgloss.Top, fainted, glow)
}

func tick(id int64, interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(_ time.Time) tea.Msg {
		return FrameMsg{id: id}
	})
}
