// Package screen provides difference screens for the tui
package screen

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/ui/view"
)

// Cammie represents the cammie screen
type Cammie struct {
	db   *db.DB
	zess *view.ZessModel
	tap  *view.TapModel
}

// NewCammie creates a new cammie screen
func NewCammie(db *db.DB) tea.Model {
	return &Cammie{db: db, zess: view.NewZessModel(db), tap: view.NewTapModel(db)}
}

// Init initializes the cammie screen
func (c *Cammie) Init() tea.Cmd {
	return tea.Batch(c.zess.Init(), c.tap.Init())
}

// Update updates the cammie screen
func (c *Cammie) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cammie, cmd := c.zess.Update(msg)
	c.zess = cammie.(*view.ZessModel)

	tap, cmd2 := c.tap.Update(msg)
	c.tap = tap.(*view.TapModel)

	return c, tea.Batch(cmd, cmd2)
}

// View returns the cammie screen view
func (c *Cammie) View() string {
	return fmt.Sprintf("%s\n%s", c.tap.View(), c.zess.View())
}
