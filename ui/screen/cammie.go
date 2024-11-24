// Package screen provides difference screens for the tui
package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/ui/view"
)

// Cammie represents the cammie screen
type Cammie struct {
	db *db.DB
}

// NewCammie creates a new cammie screen
func NewCammie(db *db.DB) Screen {
	return &Cammie{db: db}
}

// Init initializes the cammie screen
func (c *Cammie) Init() tea.Cmd {
	return nil
}

// Update updates the cammie screen
func (c *Cammie) Update(_ tea.Msg) (Screen, tea.Cmd) {
	return c, nil
}

// View returns the cammie screen view
func (c *Cammie) View() string {
	return ""
}

// GetUpdateViews returns all the update functions for the cammie screen
func (c *Cammie) GetUpdateViews() []view.UpdateData {
	return []view.UpdateData{}
}
