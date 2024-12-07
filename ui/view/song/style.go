package song

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	cZeus    = lipgloss.Color("#FF7F00")
	cSpotify = lipgloss.Color("#1DB954")
)

// Base style
var sBase = lipgloss.NewStyle()

// Styles for the stats
var (
	wStatTotal  = 30
	wStatEnum   = 3
	wStatAmount = 4
	wStatBody   = wStatTotal - wStatEnum - wStatAmount

	sStat      = sBase.Width(wStatTotal).MarginRight(3).MarginBottom(2)
	sStatTitle = sBase.Foreground(cZeus).Width(wStatTotal).Align(lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cSpotify)
	sStatEnum   = sBase.Foreground(cSpotify).Width(wStatEnum).Align(lipgloss.Left)
	sStatBody   = sBase.Width(wStatBody)
	sStatAmount = sBase.Width(wStatAmount).Align(lipgloss.Right).Foreground(cZeus)
)

// Styles for the lyrics
var (
	sLyricBase     = sBase.Width(50).Align(lipgloss.Center)
	sLyricPrevious = sLyricBase.Foreground(cZeus).Faint(true)
	sLyricCurrent  = sLyricBase.Foreground(cZeus)
	sLyricUpcoming = sLyricBase.Foreground(cSpotify)
)
