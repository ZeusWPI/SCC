package song

import "github.com/charmbracelet/lipgloss"

// Colors
var (
	cZeus    = lipgloss.Color("#FF7F00")
	cSpotify = lipgloss.Color("#1DB954")
	cBorder  = lipgloss.Color("#383838")
)

// Base style
var base = lipgloss.NewStyle()

// Styles for the stats
var (
	wStatTotal  = 40
	wStatEnum   = 3
	wStatAmount = 4
	wStatBody   = wStatTotal - wStatEnum - wStatAmount

	sStatAll   = base.Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(cBorder).PaddingTop(3)
	sStat      = base.Width(wStatTotal).MarginRight(3).MarginBottom(2)
	sStatTitle = base.Foreground(cZeus).Width(wStatTotal).Align(lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cSpotify)
	sStatEnum   = base.Foreground(cSpotify).Width(wStatEnum).Align(lipgloss.Left)
	sStatBody   = base.Width(wStatBody)
	sStatAmount = base.Width(wStatAmount).Align(lipgloss.Right).Foreground(cZeus)
)

// Styles for the status
var (
	sStatus                = base.PaddingTop(1)
	sStatusSong            = base.Padding(0, 1).Align(lipgloss.Center)
	sStatusStopwatch       = base.Faint(true)
	sStatusProgress        = base.Padding(0, 2).PaddingBottom(3).Align(lipgloss.Left)
	sStatusProgressFainted = base.Foreground(cZeus).Faint(true)
	sStatusProgressGlow    = base.Foreground(cZeus)
)

// Styles for the lyrics
var (
	sLyricBase     = base.Width(50).Align(lipgloss.Center).Bold(true)
	sLyric         = sLyricBase.AlignVertical(lipgloss.Center)
	sLyricPrevious = sLyricBase.Foreground(cZeus).Faint(true)
	sLyricCurrent  = sLyricBase.Foreground(cZeus)
	sLyricUpcoming = sLyricBase.Foreground(cSpotify)
)

// Style for everything
var (
	sAll = base.Align(lipgloss.Center).AlignVertical(lipgloss.Center)
)
