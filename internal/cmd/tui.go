// Package cmd provides all the commands to start parts of the application
package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/util"
	tui "github.com/zeusWPI/scc/ui"
	"github.com/zeusWPI/scc/ui/screen"
)

var screens = map[string]func(*db.DB) tea.Model{
	"cammie": screen.NewCammie,
}

// TUI starts the terminal user interface
func TUI(db *db.DB) (*tea.Program, error) {
	args := os.Args
	if len(args) < 2 {
		return nil, fmt.Errorf("No screen specified. Options are %v", util.Keys(screens))
	}

	screen := args[1]

	val, ok := screens[screen]
	if !ok {
		return nil, fmt.Errorf("Screen %s not found. Options are %v", screen, util.Keys(screens))
	}

	tui := tui.New(val(db))

	program := tea.NewProgram(tui)

	return program, nil
}
