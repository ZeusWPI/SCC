package gamification

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/theme"
	"github.com/zeusWPI/scc/tui/view"
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
	base.Foreground(theme.Gold),
	base.Foreground(theme.Zeus),
	base.Foreground(theme.Bronze),
	base,
}

func (m *Model) updateStyles() {
	// Adjust all style
	sAll = sAll.Width(m.width).Height(m.height).MaxHeight(m.height)

	// Adjust styles
	wColumn = (sAll.GetWidth() - view.GetOuterWidth(sAll) - view.GetOuterWidth(sColumn)*wAmount) / wAmount

	sName = sName.Width(wColumn).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(theme.Border)
	sScore = sScore.Width(wColumn)
	sAvatar = sAvatar.Width(wColumn)
}
