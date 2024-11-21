// Package screen provides difference screens for the tui
package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/ui/view"
)

// Cammie represents the cammie screen
type Cammie struct {
	db     *db.DB
	cammie *view.MessageModel
}

// NewCammie creates a new cammie screen
func NewCammie(db *db.DB) tea.Model {
	return &Cammie{db: db, cammie: view.NewMessageModel(db)}
}

// Init initializes the cammie screen
func (c *Cammie) Init() tea.Cmd {
	return c.cammie.Init()
}

// Update updates the cammie screen
func (c *Cammie) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cammie, cmd := c.cammie.Update(msg)
	c.cammie = cammie.(*view.MessageModel)

	return c, cmd
}

// View returns the cammie screen view
func (c *Cammie) View() string {
	return c.cammie.View()
}
