// Package gamification provides the functions to draw an overview of gamification on a TUI
package gamification

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/ui/view"
)

// Model represents the view model for gamification
type Model struct {
	db          *db.DB
	leaderboard []gamificationItem
}

type gamificationItem struct {
	image image.Image
	item  dto.Gamification
}

// Msg contains the data to update the gamification model
type Msg struct {
	leaderboard []gamificationItem
}

// NewModel initializes a new gamification model
func NewModel(db *db.DB) view.View {
	return &Model{db: db, leaderboard: []gamificationItem{}}
}

// Init starts the gamification view
func (m *Model) Init() tea.Cmd {
	return nil
}

// Name returns the name of the view
func (m *Model) Name() string {
	return "Gamification"
}

// Update updates the gamification view
func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		m.leaderboard = msg.leaderboard
	}

	return m, nil
}

// View draws the gamification view
func (m *Model) View() string {
	columns := make([]string, 0, len(m.leaderboard))

	positions := []lipgloss.Style{sFirst, sSecond, sThird, sFourth}

	for i, item := range m.leaderboard {
		user := lipgloss.JoinVertical(lipgloss.Left,
			positions[i].Render(fmt.Sprintf("%d. %s", i+1, item.item.Name)),
			sScore.Render(strconv.Itoa(int(item.item.Score))),
		)

		column := lipgloss.JoinVertical(lipgloss.Left, view.ImagetoString(width, item.image), user)
		columns = append(columns, sColumn.Render(column))
	}

	list := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	return list
}

// GetUpdateDatas get all update functions for the gamification view
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "gamification leaderboard",
			View:     m,
			Update:   updateLeaderboard,
			Interval: config.GetDefaultInt("tui.gamification.interval_s", 3600),
		},
	}
}

func updateLeaderboard(view view.View) (tea.Msg, error) {
	m := view.(*Model)

	gams, err := m.db.Queries.GetAllGamificationByScore(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	// Check if both leaderboards are equal
	equal := false
	if len(m.leaderboard) == len(gams) {
		equal = true
		for i, l := range m.leaderboard {
			if !l.item.Equal(*dto.GamificationDTO(gams[i])) {
				equal = false
				break
			}
		}
	}

	if equal {
		return nil, nil
	}

	msg := Msg{leaderboard: []gamificationItem{}}
	for _, gam := range gams {
		im, _, err := image.Decode(bytes.NewReader(gam.Avatar))
		if err != nil {
			return nil, err
		}

		msg.leaderboard = append(msg.leaderboard, gamificationItem{image: im, item: *dto.GamificationDTO(gam)})
	}

	return msg, nil
}
