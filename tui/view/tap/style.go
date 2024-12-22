package tap

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
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

// Base style
var base = lipgloss.NewStyle()

// Styles for the barchart
var (
	// Widths
	wBarGap    = 1                    // Gap between bars
	wBar       = 5                    // Width of bar, gets dynamically adjusted
	wBarMin    = 4                    // Minimum width required for the bar label
	wBarAmount = len(categoryToStyle) // Amount of bars. Is the same as the amount of categories

	sBar      = base.MarginBottom(1)
	sBarOne   = base
	sBarLabel = base.Align(lipgloss.Center)
)

// Styles for the stats
var (
	// Widths
	wStatAmount   = 5  // Supports up to 9999 with a space after it (or 99999 without a space)
	wStatCategory = 4  // Longest label is 4 chars
	wStatLast     = 11 // 11 characters, for example '18:53 20/12'
	wStatGapMin   = 3  // Minimum gap size between the category and last purchase

	sStat         = base.BorderStyle(lipgloss.ThickBorder()).BorderForeground(cBorder).BorderLeft(true).Margin(0, 1, 1, 1).PaddingLeft(1)
	sStatTitle    = base.Foreground(cStatsTitle).Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(cBorder).BorderBottom(true).Align(lipgloss.Center).MarginBottom(1)
	sStatAmount   = base.Width(wStatAmount).Bold(true)
	sStatCategory = base.Width(wStatCategory)
	sStatLast     = base.Width(wStatLast).Align(lipgloss.Right).Italic(true).Faint(true)
)

// Styles for the different categories
var (
	sMate = base.Foreground(cMate)
	sSoft = base.Foreground(cSoft)
	sBeer = base.Foreground(cBeer)
	sFood = base.Foreground(cFood)
)

func (m *Model) updateStyles() {
	wStatWithoutGap := wStatAmount + wStatCategory + wStatLast + view.GetOuterWidth(sStat)
	wBar = (m.width - wStatWithoutGap - wStatGapMin - wBarAmount*wBarGap) / wBarAmount
	if wBar < wBarMin {
		// Screen too small
		return
	}

	// Adjust bar styles
	sBar = sBar.Width(wBarAmount*wBar + view.GetOuterWidth(sBar)).Height(m.height - view.GetOuterHeight(sBar))
	sBarLabel = sBarLabel.Width(wBar)

	// Adjust stat styles
	wStatGap := m.width - wStatWithoutGap - (wBarAmount * wBar)
	wStat := wStatAmount + wStatCategory + wStatGap + wStatLast

	sStat = sStat.Width(wStat + view.GetOuterWidth(sStat)).Height(m.height - view.GetOuterHeight(sStat)).MaxHeight(m.height - view.GetOuterHeight(sStat))
	sStatTitle = sStatTitle.Width(wStat)
	sStatCategory = sStatCategory.Width(wStatCategory + wStatGap)
}
