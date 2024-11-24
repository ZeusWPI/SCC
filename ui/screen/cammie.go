// Package screen provides difference screens for the tui
package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/ui/view"
)

// Cammie represents the cammie screen
type Cammie struct {
	db   *db.DB
	zess view.View
}

// NewCammie creates a new cammie screen
func NewCammie(db *db.DB) Screen {
	return &Cammie{db: db, zess: view.NewZessModel(db)}
}

// Init initializes the cammie screen
func (c *Cammie) Init() tea.Cmd {
	return c.zess.Init()
}

// Update updates the cammie screen
func (c *Cammie) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case view.ZessMsg:
		zess, cmd := c.zess.Update(msg)
		c.zess = zess

		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd
	if cmds != nil {
		cmd = tea.Batch(cmds...)
	}

	return c, cmd
}

// View returns the cammie screen view
func (c *Cammie) View() string {
	return c.zess.View()
}

// GetUpdateViews returns all the update functions for the cammie screen
func (c *Cammie) GetUpdateViews() []view.UpdateData {
	return c.zess.GetUpdateDatas()
}
