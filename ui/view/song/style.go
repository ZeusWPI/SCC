package song

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	cZeus    = lipgloss.Color("#FF7F00")
	cSpotify = lipgloss.Color("#1DB954")
)

// Styles
var (
	sBase      = lipgloss.NewStyle()
	sStat      = sBase.MarginRight(3)
	sStatTitle = sBase.Foreground(cZeus).
			BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cSpotify)
	sStatAmount = sBase.Foreground(cZeus).MarginLeft(2)
	sListEnum   = sBase.Foreground(cSpotify).MarginRight(1)

	sLyricBase     = sBase.Width(50).Align(lipgloss.Center)
	sLyricPrevious = sLyricBase.Foreground(cZeus).Faint(true)
	sLyricCurrent  = sLyricBase.Foreground(cZeus)
	sLyricUpcoming = sLyricBase.Foreground(cSpotify)
)
