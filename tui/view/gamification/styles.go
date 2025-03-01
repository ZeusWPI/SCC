package gamification

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
)

// Colors
var (
	cGold   = lipgloss.Color("#FFD700")
	cZeus   = lipgloss.Color("#FF7F00")
	cBronze = lipgloss.Color("#CD7F32")
	cBorder = lipgloss.Color("#383838")
)

// Base style
var base = lipgloss.NewStyle()

// All style
var sAll = base.Align(lipgloss.Center).AlignVertical(lipgloss.Center)

// Styles
var (
	wColumn = 20 // Width of an avatar
	wAmount = 4  // Amount of people that are shown

	sColumn = base.Margin(0, 1)
	sName   = base.Align(lipgloss.Center)
	sScore  = base.Align(lipgloss.Center)
	sAvatar = base.Align(lipgloss.Center)
)

// Styles for the positions
var positions = []lipgloss.Style{
	base.Foreground(cGold),
	base.Foreground(cZeus),
	base.Foreground(cBronze),
	base,
}

func (m *Model) updateStyles() {
	// Adjust all style
	sAll = sAll.Width(m.width).Height(m.height).MaxHeight(m.height)

	// Adjust styles
	wColumn = (sAll.GetWidth() - view.GetOuterWidth(sAll) - view.GetOuterWidth(sColumn)*wAmount) / wAmount

	sName = sName.Width(wColumn).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cBorder)
	sScore = sScore.Width(wColumn)
	sAvatar = sAvatar.Width(wColumn)
}
