// Package tui provides utilities for working with the terminal.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/tui/screen"
	"go.uber.org/zap"
)

// TUI represent a terminal instance
type TUI struct {
	screen screen.Screen
}

// New creates a new tui instance
func New(screen screen.Screen) *TUI {
	return &TUI{screen: screen}
}

// Init initializes the tui
func (t *TUI) Init() tea.Cmd {
	return t.screen.Init()
}

// Update updates the tui
func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	screen, cmd := t.screen.Update(msg)
	t.screen = screen
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	// Handle global key events
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			zap.S().Info("Exiting")
			cmds = append(cmds, tea.ExitAltScreen)
			cmds = append(cmds, tea.Quit)
		}
	}

	return t, tea.Batch(cmds...)
}

// View returns the tui view
func (t *TUI) View() string {
	return t.screen.View()
}
