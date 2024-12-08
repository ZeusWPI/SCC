package gamification

import "github.com/charmbracelet/lipgloss"

var base = lipgloss.NewStyle()
var width = 20

// Colors
var (
	cGold   = lipgloss.Color("#FFD700")
	cZeus   = lipgloss.Color("#FF7F00")
	cBronze = lipgloss.Color("#CD7F32")
	cBorder = lipgloss.Color("#383838")
)

// Styles
var (
	sName   = base.BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cBorder).Width(width).Align(lipgloss.Center)
	sScore  = base.Width(width).Align(lipgloss.Center)
	sColumn = base.MarginRight(4)
)

// Positions
var (
	sFirst  = sName.Foreground(cGold)
	sSecond = sName.Foreground(cZeus)
	sThird  = sName.Foreground(cBronze)
	sFourth = sName
)
