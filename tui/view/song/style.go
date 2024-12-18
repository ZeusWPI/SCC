package song

import (
	"math"

	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
)

// Title for statistics
const (
	tStatHistory = "Recently Played"
	tStatSong    = "Top Songs"
	tStatGenre   = "Top Genres"
	tStatArtist  = "Top Artists"
)

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
	// Widths
	wStatEnum     = 3
	wStatAmount   = 4
	wStatEntryMax = 35

	// Styles
	sStat       = base.Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(cBorder).PaddingTop(1)
	sStatOne    = base.Margin(0, 1)
	sStatTitle  = base.Foreground(cZeus).Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cSpotify)
	sStatEnum   = base.Foreground(cSpotify).Width(wStatEnum).Align(lipgloss.Left)
	sStatEntry  = base.Align(lipgloss.Left)
	sStatAmount = base.Foreground(cZeus).Width(wStatAmount).Align(lipgloss.Right)
)

// Styles for the lyrics
var (
	wLyricsF = 0.8 // Fraction of width

	sLyric         = base.AlignVertical(lipgloss.Center).Align(lipgloss.Center)
	sLyricPrevious = base.Foreground(cZeus).Bold(true).Align(lipgloss.Center).Faint(true)
	sLyricCurrent  = base.Foreground(cZeus).Bold(true).Align(lipgloss.Center)
	sLyricUpcoming = base.Foreground(cSpotify).Bold(true).Align(lipgloss.Center)
)

// Styles for the status
var (
	sStatus          = base
	sStatusSong      = base.Align(lipgloss.Center)
	sStatusStopwatch = base.Faint(true)
	sStatusBar       = base.Foreground(cZeus).Align(lipgloss.Left)
)

// Style for everything
var (
	sAll = base.Align(lipgloss.Center).AlignVertical(lipgloss.Center)
)

// updateStyles updates all the affected styles when a size update message is received
func (m *Model) updateStyles() {
	// Adjust stats styles
	sStat = sStat.Width(m.width)

	wStatEntry := int(math.Min(float64(wStatEntryMax), float64(m.width/4)-float64(view.GetOuterWidth(sStatOne)+wStatEnum+wStatAmount)))
	sStatEntry = sStatEntry.Width(wStatEntry)
	sStatOne = sStatOne.Width(wStatEnum + wStatAmount + wStatEntry)
	sStatTitle = sStatTitle.Width(wStatEnum + wStatAmount + wStatEntry)
	if wStatEntry == wStatEntryMax {
		// We're full screen
		sStatOne = sStatOne.Margin(0, 3)
	}

	// Adjust lyrics styles
	sLyric = sLyric.Width(m.width)

	wLyrics := int(float64(m.width) * wLyricsF)
	sLyricPrevious = sLyricPrevious.Width(wLyrics)
	sLyricCurrent = sLyricCurrent.Width(wLyrics)
	sLyricUpcoming = sLyricUpcoming.Width(wLyrics)

	// Adjust status styles

	sStatusSong = sStatusSong.Width(m.width - view.GetOuterWidth(sStatusSong))
	sStatusBar = sStatusBar.Width(m.width - view.GetOuterWidth(sStatusBar))

	// Adjust the all styles
	sAll = sAll.Height(m.height - view.GetOuterHeight(sAll)).Width(m.width - view.GetOuterWidth(sAll))
}
