// Package screen provides difference screens for the tui
package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/ui/view"
)

// Screen represents a screen
type Screen interface {
	Init() tea.Cmd
	Update(tea.Msg) (Screen, tea.Cmd)
	View() string
	GetUpdateViews() []view.UpdateData
}
