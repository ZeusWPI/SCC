package message

import "github.com/charmbracelet/lipgloss"

var base = lipgloss.NewStyle()

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

// Style
var (
	sTime    = base.Faint(true)
	sSender  = base.Bold(true)
	sMessage = base
	sDate    = base.Faint(true)
)
