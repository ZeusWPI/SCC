// Package cmd provides all the commands to start parts of the application
package cmd

import (
	"context"
	"fmt"
	"maps"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui"
	"github.com/zeusWPI/scc/tui/screen"
	"github.com/zeusWPI/scc/tui/screen/cammie"
	songScreen "github.com/zeusWPI/scc/tui/screen/song"
	"github.com/zeusWPI/scc/tui/view"
)

var screens = map[string]func(repo repository.Repository) screen.Screen{
	"cammie": cammie.New,
	"song":   songScreen.New,
}

func TUI(repo repository.Repository, screenName string) error {
	val, ok := screens[screenName]
	if !ok {
		return fmt.Errorf("screen %s not found. Options are %v", screenName, maps.Keys(screens))
	}

	screen := val(repo)
	tui := tui.New(screen)
	p := tea.NewProgram(tui, tea.WithAltScreen())

	dones := make([]chan bool, 0, len(screen.GetUpdateViews()))
	for _, data := range screen.GetUpdateViews() {
		dones = append(dones, periodicUpdate(p, data))
	}

	_, err := p.Run()

	for _, done := range dones {
		done <- true
	}

	return err
}

func periodicUpdate(p *tea.Program, data view.UpdateData) chan bool {
	done := make(chan bool)

	update := func(ctx context.Context) error {
		msg, err := data.Update(ctx, data.View)
		if err != nil {
			return err
		}

		if msg != nil {
			p.Send(msg)
		}

		return nil
	}

	go utils.Periodic(
		data.Name,
		time.Duration(data.Interval)*time.Second,
		update,
		done,
	)

	return done
}
