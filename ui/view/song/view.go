package song

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewPlaying() string {
	var previousB strings.Builder
	for i, lyric := range m.current.previous {
		previousB.WriteString(lyric)
		if i != len(m.current.previous)-1 {
			previousB.WriteString("\n")
		}
	}
	previous := sLyricPrevious.Render(previousB.String())

	current := sLyricCurrent.Render(m.current.current)

	var upcomingB strings.Builder
	for _, lyric := range m.current.upcoming {
		upcomingB.WriteString(lyric)
		upcomingB.WriteString("\n")
	}
	upcoming := sLyricUpcoming.Render(upcomingB.String())

	return sBase.MarginLeft(5).Render(lipgloss.JoinVertical(lipgloss.Left, previous, current, upcoming))
}

func (m *Model) viewNotPlaying() string {
	rows := make([][]string, 0, 2)
	for i := 0; i < 2; i++ {
		rows = append(rows, make([]string, 0, 2))
	}

	// Recently played
	items := make([]string, 0, len(m.history))
	for i, track := range m.history {
		number := sStatEnum.Render(fmt.Sprintf("%d.", i+1))
		body := sStatBody.Render(track)
		items = append(items, lipgloss.JoinHorizontal(lipgloss.Top, number, body))
	}
	l := lipgloss.JoinVertical(lipgloss.Left, items...)
	title := sStatTitle.Render("Recently Played")
	rows[0] = append(rows[0], sStat.Render(lipgloss.JoinVertical(lipgloss.Left, title, l)))

	// All other stats
	topStats := [][]topStat{m.topSongs, m.topArtists, m.topGenres}
	for i, topStat := range topStats {
		items := make([]string, 0, len(topStat))
		for i, stat := range topStat {
			number := sStatEnum.Render(fmt.Sprintf("%d.", i+1))
			body := sStatBody.Render(stat.name)
			amount := sStatAmount.Render(fmt.Sprintf("%d", stat.amount))
			items = append(items, lipgloss.JoinHorizontal(lipgloss.Top, number, body, amount))
		}
		l := lipgloss.JoinVertical(lipgloss.Left, items...)

		var row int
		if i == 0 {
			title = sStatTitle.Render("Top Tracks")
			row = 0
		} else if i == 1 {
			title = sStatTitle.Render("Top Artists")
			row = 1
		} else {
			title = sStatTitle.Render("Top Genres")
			row = 1
		}

		rows[row] = append(rows[row], sStat.Render(lipgloss.JoinVertical(lipgloss.Left, title, l)))
	}

	renderedRows := make([]string, 0, 2)
	for _, row := range rows {
		renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, renderedRows...)
}
