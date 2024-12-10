package tap

import (
	"strconv"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewChart() string {
	chart := barchart.New(widthBar, heightBar)
	bars := make([]barchart.BarData, 0, len(m.items))

	for _, item := range m.items {
		style, ok := categoryToStyle[item.category]
		if !ok {
			style = sUnknown
		}

		bars = append(bars, barchart.BarData{
			Label: string(item.category),
			Values: []barchart.BarValue{{
				Name:  string(item.category),
				Value: float64(item.amount),
				Style: style,
			}},
		})
	}

	chart.PushAll(bars)
	chart.Draw()

	return chart.View()
}

func (m *Model) viewStats() string {
	rows := make([]string, 0, len(m.items))

	for _, item := range m.items {
		amount := sStatsAmount.Render(strconv.Itoa(item.amount))
		category := sStatsCategory.Inherit(categoryToStyle[item.category]).Render(string(item.category))
		last := sStatsLast.Render(item.last.Format("15:04 02/01"))

		text := lipgloss.JoinHorizontal(lipgloss.Top, amount, category, last)
		rows = append(rows, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Add title
	title := sStatsTitle.Render("Leaderboard")
	view = lipgloss.JoinVertical(lipgloss.Left, title, view)

	return view
}
