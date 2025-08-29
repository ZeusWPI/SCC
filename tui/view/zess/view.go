package zess

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/pkg/utils"
)

func (m *Model) viewChart() string {
	chart := barchart.New(sBar.GetWidth(), sBar.GetHeight(), barchart.WithNoAutoBarWidth(), barchart.WithBarGap(wBarGap), barchart.WithBarWidth(wBar))

	for _, week := range m.weeks {
		weekNumber, _ := week.start.ISOWeek()

		bar := barchart.BarData{
			Label: sBarLabel.Render(fmt.Sprintf("W%02d", weekNumber)),
			Values: []barchart.BarValue{{
				Name:  week.start.String(),
				Value: float64(week.scans),
				Style: sBarOne.Foreground(lipgloss.Color(randomColor())),
			}},
		}

		chart.Push(bar)
	}

	chart.Draw()

	return chart.View()
}

func (m *Model) viewStats() string {
	// Overview of each week
	rows := make([]string, 0, len(m.weeks))

	maxScans := slices.MaxFunc(m.weeks, func(a, b week) int { return a.scans - b.scans })

	for _, week := range m.weeks {
		weekNumber, _ := week.start.ISOWeek()
		weekStr := sStatDate.Render(fmt.Sprintf("W%02d - %s", weekNumber, week.start.Format("01/02")))

		var amount string
		if week.scans == maxScans.scans {
			amount = sStatMax.Inherit(sStatAmount).Render(strconv.Itoa(week.scans))
		} else {
			amount = sStatAmount.Render(strconv.Itoa(week.scans))
		}

		text := lipgloss.JoinHorizontal(lipgloss.Top, weekStr, amount)
		rows = append(rows, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Title
	title := sStatTitle.Render("Overview")

	// Total scans
	total := sStatTotalTitle.Render("Total")
	amount := sStatTotalAmount.Render(strconv.Itoa(utils.Reduce(m.weeks, func(accum int, w week) int { return accum + w.scans })))
	total = lipgloss.JoinHorizontal(lipgloss.Top, total, amount)
	total = sStatTotal.Render(total)

	view = lipgloss.JoinVertical(lipgloss.Left, title, view, total)

	return sStat.Render(view)
}
