// Package view contains all the different views for the tui
package view

import (
	tea "github.com/charmbracelet/bubbletea"
)

// UpdateData represents the data needed to update a view
type UpdateData struct {
	Name     string
	View     View
	Update   func(view View) (tea.Msg, error)
	Interval int
}

// View represents a view
type View interface {
	Init() tea.Cmd
	Update(tea.Msg) (View, tea.Cmd)
	View() string
	GetUpdateDatas() []UpdateData
}
