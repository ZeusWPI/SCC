package zess

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
)

// Colors
var (
	cBorder     = lipgloss.Color("#383838")
	cZeus       = lipgloss.Color("#FF7F00")
	cStatsTitle = lipgloss.Color("#EE4B2B")
)

// Base style
var base = lipgloss.NewStyle()

// Styles for the barchart
var (
	// Widths
	wBarGap       = 1  // Gap between bars
	wBar          = 5  // Width of a bar
	wBarMin       = 3  // Required for the bar label, for example 'W56'
	wBarAmountMax = 10 // Maximum amount of bars

	sBar      = base
	sBarOne   = base
	sBarLabel = base.Align(lipgloss.Center)
)

// Styles for the stats
var (
	// Widths
	wStatDate   = 11 // 11 characters, for example 'W56 - 29/12'
	wStatAmount = 4  // Supports up to 9999
	wStatGapMin = 3  // Minumum gap size between the date and amount

	sStat            = base.BorderStyle(lipgloss.ThickBorder()).BorderForeground(cBorder).BorderLeft(true).MarginLeft(1).PaddingLeft(1)
	sStatTitle       = base.Foreground(cStatsTitle).Bold(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(cBorder).BorderBottom(true).Align(lipgloss.Center).MarginBottom(1)
	sStatDate        = base.Width(wStatDate)
	sStatAmount      = base.Width(wStatAmount)
	sStatTotal       = base.BorderStyle(lipgloss.NormalBorder()).BorderForeground(cBorder).BorderTop(true).MarginTop(1)
	sStatTotalTitle  = sStatDate.Bold(true)
	sStatTotalAmount = sStatAmount.Bold(true)
)

// Styles for the max amount
var (
	sMax = base.Foreground(cZeus).Bold(true)
)

// Bar colors
var colors = []string{
	"#FAF500", // Yellow
	"#3AFA00", // Green
	"#FAD700", // Yellow Green
	"#FAA600", // Orange
	"#FAE200", // Yellow Orange
	"#FA7200", // Orange Red
	"#FA4600", // Red
	"#FA0400", // Real Red
	"#FA0079", // Pink Red
	"#FA00FA", // Pink
	"#EE00FA", // Purple
	"#8300FA", // Purple Blue
	"#3100FA", // Blue
	"#00FAFA", // Light Blue
	"#00FAA5", // Green Blue
	"#00FA81", // IDK
	"#F8FA91", // Weird Light Green
	"#FAD392", // Light Orange
	"#FA9E96", // Salmon
	"#DEA2F9", // Fuchsia
	"#B3D2F9", // Boring Blue
}

// updateStyles updates all the affected styles when a size update message is received
func (m *Model) updateStyles() {
	if m.width-wStatDate-wStatAmount-wStatGapMin < 0 {
		// Screen is way too small
		// Avoid entering an infinite loop down below
		return
	}

	// Adjust bar styles

	wBar = wBarMin
	wStatWithoutGap := wStatDate + wStatAmount + view.GetOuterWidth(sStat)
	for (m.width-wStatWithoutGap-wStatGapMin-wBarAmountMax*wBarGap)/wBar >= wBarAmountMax {
		wBar++
	}
	bars := (m.width - wStatWithoutGap - wStatGapMin) / wBar
	sBar = sBar.Width(bars * wBar).Height(m.height - view.GetOuterHeight(sBar))
	sBarLabel = sBarLabel.Width(wBar)

	// Adjust stat styles
	wStatGap := m.width - wStatWithoutGap - (bars * wBar)
	wStat := wStatDate + wStatGap + wStatAmount

	sStat = sStat.Width(wStat + view.GetOuterWidth(sStat)).Height(m.height).MaxHeight(m.height)
	sStatTitle = sStatTitle.Width(wStat)
	sStatDate = sStatDate.Width(sStatDate.GetWidth() + wStatGap)
	sStatTotal = sStatTotal.Width(sStatTitle.GetWidth())
	sStatTotalTitle = sStatTotalTitle.Width(sStatDate.GetWidth())

}

// // Width
// var (
// 	widthAmount = 5
// 	widthWeek   = 8
// )

// // Margin
// var mOverview = 2

// // Barchart
// var (
// 	widthBar  = 60
// 	heightBar = 20
// )

// // Styles stats
// var (
// 	sStats            = base.Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(cBorder).MarginLeft(mOverview).PaddingLeft(mOverview)
// 	sStatsTitle       = base.Foreground(cStatsTitle).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(cBorder).Width(widthAmount + widthWeek).Align(lipgloss.Center)
// 	sStatsWeek        = base.Width(widthWeek)
// 	sStatsAmount      = base.Bold(true).Width(widthAmount).Align(lipgloss.Right)
// 	sStatsAmountMax   = sStatsAmount.Foreground(cZeus)
// 	sStatsTotal       = base.Border(lipgloss.NormalBorder(), true, false, false, false).BorderForeground(cBorder).MarginTop(1)
// 	sStatsTotalTitle  = sStatsWeek
// 	sStatsTotalAmount = sStatsAmount
// )
