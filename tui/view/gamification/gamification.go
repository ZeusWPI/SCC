// Package gamification provides the functions to draw an overview of gamification on a TUI
package gamification

import (
	"fmt"
	"image"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/view"
)

type Model struct {
	leaderboard []gamification

	width  int
	height int

	url string // API url for gamification
}

// Interface compliance
var _ view.View = (*Model)(nil)

// Msg contains the data to update the gamification model
type Msg struct {
	leaderboard []gamification
}

// Interface Compliance
var _ tea.Msg = (*Msg)(nil)

type gamification struct {
	Name      string `json:"github_name"`
	Score     int    `json:"score"`
	AvatarURL string `json:"avartar_url"`
	avatar    image.Image
}

func NewModel() view.View {
	return &Model{
		leaderboard: nil,
		width:       0,
		height:      0,
		url:         config.GetDefaultString("tui.view.gamification.url", "https://gamification.zeus.gent"),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Name() string {
	return "Gamification"
}

func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		entry, ok := msg.Sizes[m.Name()]
		if ok {
			// Update all dependent styles
			m.width = entry.Width
			m.height = entry.Height

			m.updateStyles()
		}

		return m, nil

	case Msg:
		m.leaderboard = msg.leaderboard
	}

	return m, nil
}

func (m *Model) View() string {
	columns := make([]string, 0, len(m.leaderboard))

	for i, item := range m.leaderboard {
		user := lipgloss.JoinVertical(lipgloss.Left,
			positions[i].Inherit(sName).Render(fmt.Sprintf("%d. %s", i+1, item.Name)),
			sScore.Render(strconv.Itoa(int(item.Score))),
		)
		im := sAvatar.Render(view.ImageToString(item.avatar, wColumn, sAll.GetHeight()-lipgloss.Height(user)))

		column := lipgloss.JoinVertical(lipgloss.Left, im, user)
		columns = append(columns, sColumn.Render(column))
	}

	list := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	return sAll.Render(list)
}

// GetUpdateDatas get all update functions for the gamification view
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "gamification leaderboard",
			View:     m,
			Update:   updateLeaderboard,
			Interval: config.GetDefaultInt("tui.view.gamification.interval_s", 3600),
		},
	}
}
