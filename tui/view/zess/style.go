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
	cBorder     = lipgloss.Color("#383838")
	cZeus       = lipgloss.Color("#FF7F00")
	cStatsTitle = lipgloss.Color("#EE4B2B")
)

// Message colors
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

// Styles chart
var (
	sBar = base
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
