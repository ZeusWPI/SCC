// Package screen provides difference screens for the tui
package screen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
)

// Cammie represents the cammie screen
type Cammie struct {
	db *db.DB
}

// NewCammie creates a new cammie screen
func NewCammie(db *db.DB) tea.Model {
	return &Cammie{db: db}
}

// Init initializes the cammie screen
func (c *Cammie) Init() tea.Cmd {
	return nil
}

// Update updates the cammie screen
func (c *Cammie) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

// View returns the cammie screen view
func (c *Cammie) View() string {
	return ""
}
