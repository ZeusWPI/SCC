package song

import "github.com/charmbracelet/lipgloss"

var base = lipgloss.NewStyle()

// Style
var (
	sSong = base.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center)
)
