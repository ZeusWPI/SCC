// Package tui provides utilities for working with the terminal.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/ui/views"
	"go.uber.org/zap"
)

// TUI represent a terminal instance
type TUI struct {
	db  *db.DB
	tap tea.Model
}

// New creates a new tty instance
func New(db *db.DB) *TUI {
	return &TUI{
		db:  db,
		tap: views.NewTapModel(db),
	}
}

// Init initializes the tty
func (t *TUI) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, t.tap.Init())
}

// Update updates the tty
func (t *TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	tapModel, tapCmd := t.tap.Update(msg)
	if tapCmd != nil {
		cmds = append(cmds, tapCmd)
	}
	t.tap = tapModel

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

// View returns the tty view
func (t *TUI) View() string {
	return t.tap.View()
}
