// Package tui provides utilities for working with the terminal.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

// TUI represent a terminal instance
type TUI struct {
	screen tea.Model
}

// New creates a new tui instance
func New(screen tea.Model) *TUI {
	return &TUI{screen: screen}
}

// Init initializes the tui
func (t *TUI) Init() tea.Cmd {
	return tea.Batch(t.screen.Init())
}

// Update updates the tui
func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	model, cmd := t.screen.Update(msg)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	t.screen = model

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

// View returns the ttuity view
func (t *TUI) View() string {
	return t.screen.View()
}
