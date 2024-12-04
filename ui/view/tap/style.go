package tap

import "github.com/charmbracelet/lipgloss"

var base = lipgloss.NewStyle()

// Width
var (
	widthAmount   = 5
	widthCategory = 8
	widthLast     = 13
)

// Margin
var mStats = 2

// Barchart
var (
	widthBar  = 40
	heightBar = 20
)

// Colors
var (
	cMate = lipgloss.Color("#D27D2D")
	cSoft = lipgloss.Color("#ADD8E6")
	cBeer = lipgloss.Color("#F9B116")
	cFood = lipgloss.Color("#00ff00")

	cBorder     = lipgloss.Color("#383838")
	cStatsTitle = lipgloss.Color("#EE4B2B")
)

// Styles Chart
var (
	sMate    = base.Foreground(cMate)
	sSoft    = base.Foreground(cSoft)
	sBeer    = base.Foreground(cBeer)
	sFood    = base.Foreground(cFood)
	sUnknown = base
)

// Styles stats
var (
	sStats         = base.Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(cBorder).MarginLeft(mStats).PaddingLeft(mStats)
	sStatsTitle    = base.Foreground(cStatsTitle).Bold(true).Width(widthAmount+widthCategory+widthLast).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(cBorder)
	sStatsAmount   = base.Width(widthAmount).Bold(true)
	sStatsCategory = base.Width(widthCategory)
	sStatsLast     = base.Width(widthLast).Align(lipgloss.Right).Italic(true).Faint(true)
)
