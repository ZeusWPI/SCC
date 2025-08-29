// Package song provides the functions to draw an overview of the song integration
package song

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/tui/view"
)

type Model struct {
	width  int
	height int
}

// Interface compliance
var _ view.View = (*Model)(nil)

// Msg contains the data to update the gamification model
type Msg struct{}

// New initializes a new song model
func New() view.View {
	return &Model{
		width:  0,
		height: 0,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Name() string {
	return "Song"
}

func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		if entry, ok := msg.Sizes[m.Name()]; ok {
			// Update all dependent styles
			m.width = entry.Width
			m.height = entry.Height
		}

		return m, nil

	default:
		break
	}

	return m, nil
}

// View draws the song view
func (m *Model) View() string {
	return "Not implemented"
}

// GetUpdateDatas gets all update functions for the song view
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return nil
}
