package event

import "github.com/charmbracelet/lipgloss"

// Widths
var (
	widthToday = 45
	widthImage = 32

	widthOverview      = 45
	widthOverviewName  = 35
	widthOverviewImage = 32
)

// Base
var (
	base      = lipgloss.NewStyle()
	baseToday = base.Width(widthToday).Align(lipgloss.Center)
)

// Margins
var (
	mTodayWarning = 3
	mOverview     = 5
)

// Color
var (
	cZeus     = lipgloss.Color("#FF7F00")
	cWarning  = lipgloss.Color("#EE4B2B")
	cBorder   = lipgloss.Color("#383838")
	cUpcoming = lipgloss.Color("#FFBF00")
)

// Styles today
var (
	sTodayWarning = baseToday.Bold(true).Blink(true).Foreground(cWarning).Border(lipgloss.DoubleBorder(), true, false)
	sTodayName    = baseToday.Bold(true).Foreground(cZeus).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cBorder)
	sTodayTime    = baseToday
	sTodayPlace   = baseToday.Italic(true).Faint(true)
	sToday        = baseToday.MarginLeft(8).AlignVertical(lipgloss.Center)
)

// Styles overview
var (
	sOverviewTitle = base.Bold(true).Foreground(cWarning).Width(widthOverview).Align(lipgloss.Center)
	sOverview      = base.Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(cBorder).Width(widthOverview).MarginRight(mOverview)
	sPassedName    = base.Foreground(cZeus).Faint(true).Width(widthOverviewName)
	sPassedTime    = base.Faint(true)
	sNextName      = base.Bold(true).Foreground(cZeus).Width(widthOverviewName)
	sNextTime      = base.Bold(true)
	sNextPlace     = base.Italic(true).Width(widthOverviewName)
	sUpcomingName  = base.Width(widthOverviewName).Foreground(cUpcoming)
	sUpcomingTime  = base.Faint(true)
	sUpcomingPlace = base.Italic(true).Faint(true).Width(widthOverviewName)
)
