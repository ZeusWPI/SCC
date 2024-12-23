package event

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
)

// Color
var (
	cZeus     = lipgloss.Color("#FF7F00")
	cWarning  = lipgloss.Color("#EE4B2B")
	cBorder   = lipgloss.Color("#383838")
	cUpcoming = lipgloss.Color("#FFBF00")
)

// Base style
var base = lipgloss.NewStyle()

// Styles for overview
var (
	wOvDate    = 8  // Width of the date, for example '13/11' (with some padding after)
	wOvTextMin = 20 // Minimum width of the event name
	wOvPoster  = 20 // Width of the poster
	wOvGap     = 2  // Width of the gap between the overview and the poster

	sOvAll    = base.Padding(0, 1) // Style for the overview and the poster
	sOvPoster = base.AlignVertical(lipgloss.Center)
	sOv       = base.AlignVertical(lipgloss.Center).MarginRight(wOvGap) // Style for the overview of the events
	sOvTitle  = base.Bold(true).Foreground(cWarning).Align(lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cBorder)

	// Styles for passed events
	sOvPassedDate = base.Width(wOvDate).Faint(true)
	sOvPassedText = base.Foreground(cZeus).Faint(true)

	// Styles for next event
	sOvNextDate = base.Width(wOvDate).Bold(true)
	sOvNextText = base.Bold(true).Foreground(cZeus)
	sOvNextLoc  = base.Italic(true)

	// Styles for the upcoming envets
	sOvUpcomingDate = base.Width(wOvDate).Faint(true)
	sOvUpcomingText = base.Foreground(cUpcoming)
	sOvUpcomingLoc  = base.Italic(true).Faint(true)
)

// Styles for today
var (
	wTodayEvMin  = 20 // Minimum width of the event
	wTodayPoster = 20 // Width of the poster
	wTodayGap    = 2  // Width of the gap between the text and the poster

	sTodayAll    = base.Padding(0, 1) // Style for the text and the poster
	sTodayPoster = base.AlignVertical(lipgloss.Center)
	sToday       = base.AlignVertical(lipgloss.Center).MarginLeft(wOvGap).Padding(1, 0).Border(lipgloss.DoubleBorder(), true, false) // Style for the event

	sTodayDate = base.Align(lipgloss.Center)
	sTodayText = base.Align(lipgloss.Center).Bold(true).Foreground(cZeus).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(cBorder)
	sTodayeLoc = base.Align(lipgloss.Center).Italic(true).Faint(true)
)

func (m *Model) updateStyles() {
	// Adjust the styles for the overview
	wOvPoster = (m.width - wOvGap - view.GetOuterWidth(sOvAll)) / 2
	if wOvPoster <= wOvDate+wOvTextMin {
		// Screen is too small, don't draw the poster for more space
		wOvPoster = 0
	}

	wOv := wOvPoster
	wOvText := wOv - wOvDate

	sOv = sOv.Width(wOv)
	sOvTitle = sOvTitle.Width(wOv)
	sOvPassedText = sOvPassedText.Width(wOvText)
	sOvNextText = sOvNextText.Width(wOvText)
	sOvNextLoc = sOvNextLoc.Width(wOvText)
	sOvUpcomingText = sOvUpcomingText.Width(wOvText)
	sOvUpcomingLoc = sOvUpcomingLoc.Width(wOvText)

	// Adjust the styles for today
	wTodayPoster = (m.width - wTodayGap - view.GetOuterWidth(sTodayAll)) / 2
	if wTodayPoster <= wTodayEvMin {
		// Screen is too small, don't draw the poster for more space
		wTodayPoster = 0
	}

	wTodayEv := wTodayPoster

	sTodayDate = sTodayDate.Width(wTodayEv)
	sTodayText = sTodayText.Width(wTodayEv)
	sTodayeLoc = sTodayeLoc.Width(wTodayEv)
}
