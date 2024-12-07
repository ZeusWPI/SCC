package cammie

import "github.com/charmbracelet/lipgloss"

// Base
var base = lipgloss.NewStyle()

// Colors
var (
	cZeus = lipgloss.Color("#FF7F00")
)

// Borders
var (
	bTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	bActiveTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}
)

// Tabs
var (
	sTab       = base.Border(bTabBorder, true).BorderForeground(cZeus).Padding(0, 1).MarginBottom(1)
	sActiveTab = sTab.Border(bActiveTabBorder, true)
	sTabNormal = sTab.BorderTop(false).BorderLeft(false).BorderRight(false)
)

// Style
var (
	sMsg    = base.Border(lipgloss.RoundedBorder(), true, true, true, true).BorderForeground(cZeus).MarginLeft(1).MarginRight(2)
	sTop    = base.Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(cZeus)
	sBottom = base.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center)
)
