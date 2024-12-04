package zess

import (
	"strconv"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewChart() string {
	chart := barchart.New(widthBar, heightBar)

	for _, scan := range m.scans {
		bar := barchart.BarData{
			Label: scan.label,
			Values: []barchart.BarValue{{
				Name:  scan.label,
				Value: float64(scan.amount),
				Style: sBar,
			}},
		}

		chart.Push(bar)
	}

	chart.Draw()

	return chart.View()
}

func (m *Model) viewStats() string {
	// Overview of each week
	rows := make([]string, 0, len(m.scans))

	for _, scan := range m.scans {
		week := sStatsWeek.Render(scan.label)

		var amount string
		if scan.amount == m.maxWeekScans {
			amount = sStatsAmountMax.Render(strconv.Itoa(int(scan.amount)))
		} else {
			amount = sStatsAmount.Render(strconv.Itoa(int(scan.amount)))
		}

		text := lipgloss.JoinHorizontal(lipgloss.Top, week, amount)
		rows = append(rows, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Title
	title := sStatsTitle.Render("Overview")

	// Total scans
	total := sStatsTotalTitle.Render("Total")
	amount := sStatsTotalAmount.Render(strconv.Itoa(int(m.seasonScans)))
	total = lipgloss.JoinHorizontal(lipgloss.Top, total, amount)
	total = sStatsTotal.Render(total)

	view = lipgloss.JoinVertical(lipgloss.Left, title, view, total)

	return view
}
