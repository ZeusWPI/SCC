package zess

import "github.com/charmbracelet/lipgloss"

var base = lipgloss.NewStyle()

// Width
var (
	widthAmount = 5
	widthWeek   = 8
)

// Margin
var mOverview = 2

// Barchart
var (
	widthBar  = 60
	heightBar = 20
)

// Colors
var (
	cBarChart = lipgloss.Color("#32012F")

	cBorder     = lipgloss.Color("#383838")
	cZeus       = lipgloss.Color("#FF7F00")
	cStatsTitle = lipgloss.Color("#EE4B2B")
)

// Styles chart
var (
	sBar = base.Foreground(cBarChart)
)

// Styles stats
var (
	sStats            = base.Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(cBorder).MarginLeft(mOverview).PaddingLeft(mOverview)
	sStatsTitle       = base.Foreground(cStatsTitle).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(cBorder).Width(widthAmount + widthWeek).Align(lipgloss.Center)
	sStatsWeek        = base.Width(widthWeek)
	sStatsAmount      = base.Bold(true).Width(widthAmount).Align(lipgloss.Right)
	sStatsAmountMax   = sStatsAmount.Foreground(cZeus)
	sStatsTotal       = base.Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(cBorder).MarginTop(1)
	sStatsTotalTitle  = sStatsWeek
	sStatsTotalAmount = sStatsAmount
)
