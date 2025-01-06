package zess

import (
	"fmt"
	"strconv"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewChart() string {
	chart := barchart.New(sBar.GetWidth(), sBar.GetHeight(), barchart.WithNoAutoBarWidth(), barchart.WithBarGap(wBarGap), barchart.WithBarWidth(wBar))

	for _, scan := range m.scans {
		bar := barchart.BarData{
			Label: sBarLabel.Render(fmt.Sprintf("W%d", scan.time.week)),
			Values: []barchart.BarValue{{
				Name:  scan.start,
				Value: float64(scan.amount),
				Style: sBarOne.Foreground(lipgloss.Color(scan.color)),
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
		week := sStatDate.Render(fmt.Sprintf("W%02d - %s", scan.time.week, scan.start))

		var amount string
		if scan.amount == m.maxWeekScans {
			amount = sStatMax.Inherit(sStatAmount).Render(strconv.Itoa(int(scan.amount)))
		} else {
			amount = sStatAmount.Render(strconv.Itoa(int(scan.amount)))
		}

		text := lipgloss.JoinHorizontal(lipgloss.Top, week, amount)
		rows = append(rows, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Title
	title := sStatTitle.Render("Overview")

	// Total scans
	total := sStatTotalTitle.Render("Total")
	amount := sStatTotalAmount.Render(strconv.Itoa(int(m.seasonScans)))
	total = lipgloss.JoinHorizontal(lipgloss.Top, total, amount)
	total = sStatTotal.Render(total)

	view = lipgloss.JoinVertical(lipgloss.Left, title, view, total)

	return sStat.Render(view)
}
