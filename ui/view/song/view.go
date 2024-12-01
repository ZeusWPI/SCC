package song

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
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
	columns := make([]string, 0, 3)

	// Recently played
	l := list.New(m.history).Enumerator(list.Arabic).EnumeratorStyle(sListEnum).String()
	t := sStatTitle.Width(lipgloss.Width(l)).Align(lipgloss.Center).Render("Recently played")

	column := lipgloss.JoinVertical(lipgloss.Left, t, l)
	columns = append(columns, sStat.Render(column))

	// Top stats
	topStats := map[string][]topStat{
		"Top Tracks":  m.topSongs,
		"Top Artists": m.topArtists,
		"Top Genres":  m.topGenres,
	}

	for title, stat := range topStats {
		var statInfos []string
		for _, statInfo := range stat {
			statInfos = append(statInfos, lipgloss.JoinHorizontal(lipgloss.Top, statInfo.name, sStatAmount.Render(fmt.Sprintf("%d", statInfo.amount))))
		}
		l := list.New(statInfos).Enumerator(list.Arabic).EnumeratorStyle(sListEnum).String()
		t := sStatTitle.Width(lipgloss.Width(l)).Align(lipgloss.Center).Render(title)

		column := lipgloss.JoinVertical(lipgloss.Left, t, l)
		columns = append(columns, sStat.Render(column))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}
