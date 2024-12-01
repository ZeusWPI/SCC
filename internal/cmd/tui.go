// Package cmd provides all the commands to start parts of the application
package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/util"
	tui "github.com/zeusWPI/scc/ui"
	"github.com/zeusWPI/scc/ui/screen"
	"github.com/zeusWPI/scc/ui/view"
	"go.uber.org/zap"
)

var screens = map[string]func(*db.DB) screen.Screen{
	"cammie": screen.NewCammie,
	"song":   screen.NewSong,
	"test":   screen.NewTest,
}

// TUI starts the terminal user interface
func TUI(db *db.DB) error {
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("No screen specified. Options are %v", util.Keys(screens))
	}

	selectedScreen := args[1]

	val, ok := screens[selectedScreen]
	if !ok {
		return fmt.Errorf("Screen %s not found. Options are %v", selectedScreen, util.Keys(screens))
	}

	screen := val(db)
	tui := tui.New(screen)
	p := tea.NewProgram(tui, tea.WithAltScreen())

	dones := make([]chan bool, 0, len(screen.GetUpdateViews()))
	for _, updateData := range screen.GetUpdateViews() {
		done := make(chan bool)
		dones = append(dones, done)
		go tuiPeriodicUpdates(p, updateData, done)
	}

	_, err := p.Run()

	for _, done := range dones {
		done <- true
	}

	return err
}

func tuiPeriodicUpdates(p *tea.Program, updateData view.UpdateData, done chan bool) {
	zap.S().Info("TUI: Starting periodic update for ", updateData.Name, " with an interval of ", updateData.Interval, " seconds")

	ticker := time.NewTicker(time.Duration(updateData.Interval) * time.Second)
	defer ticker.Stop()

	// Immediatly update once
	msg, err := updateData.Update(updateData.View)
	if err != nil {
		zap.S().Error("TUI: Error updating ", updateData.Name, "\n", err)
	}

	if msg != nil {
		p.Send(msg)
	}

	for {
		select {
		case <-done:
			zap.S().Info("TUI: Stopping periodic update for ", updateData.Name)
			return
		case <-ticker.C:
			// Update
			msg, err := updateData.Update(updateData.View)
			if err != nil {
				zap.S().Error("TUI: Error updating ", updateData.Name, "\n", err)
			}

			if msg != nil {
				p.Send(msg)
			}
		}
	}
}
