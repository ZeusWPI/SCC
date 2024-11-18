// Package cmd provides all the commands to start parts of the application
package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	tui "github.com/zeusWPI/scc/ui"
)

// TUI starts the terminal user interface
func TUI(db *db.DB) *tea.Program {
	tui := tui.New(db)

	program := tea.NewProgram(tui)

	return program
}
