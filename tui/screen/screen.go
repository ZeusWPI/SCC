// Package screen contains the interface for a screen
package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/tui/view"
)

type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg) (Screen, tea.Cmd)
	View() string
	GetUpdateViews() []view.UpdateData
	GetSizeMsg() tea.Msg
}
