package song

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewPlaying() string {
	status := m.viewPlayingStatus()
	status = sStatus.Render(status)

	stats := m.viewPlayingStats()
	stats = sStatAll.Render(stats)

	lyrics := m.viewPlayingLyrics()
	lyrics = sLyric.Height(sAll.GetHeight() - lipgloss.Height(status) - lipgloss.Height(stats)).Render(lyrics)

	view := lipgloss.JoinVertical(lipgloss.Left, status, lyrics, stats)

	return sAll.Render(view)
}

func (m *Model) viewPlayingStatus() string {
	// Stopwatch
	durationS := int(math.Round(float64(m.current.song.DurationMS) / 1000))
	stopwatch := fmt.Sprintf("\t%s / %02d:%02d", m.current.stopwatch.View(), durationS/60, durationS%60)
	stopwatch = sStatusStopwatch.Render(stopwatch)

	// Song name
	var artists strings.Builder
	for _, artist := range m.current.song.Artists {
		artists.WriteString(artist.Name + " & ")
	}
	artist := artists.String()
	if len(artist) > 0 {
		artist = artist[:len(artist)-3]
	}

	song := sStatusSong.Width(sStatusSong.GetWidth() - lipgloss.Width(stopwatch)).Render(fmt.Sprintf("%s | %s", m.current.song.Title, artist))

	// Progress bar
	progress := m.current.progress.View()
	progress = sStatusProgress.Render(progress)

	view := lipgloss.JoinHorizontal(lipgloss.Top, song, stopwatch)
	view = lipgloss.JoinVertical(lipgloss.Left, view, progress)

	return view
}

func (m *Model) viewPlayingLyrics() string {
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

	return sLyric.Render(lipgloss.JoinVertical(lipgloss.Left, previous, current, upcoming))
}

func (m *Model) viewPlayingStats() string {
	columns := make([]string, 0, 4)

	columns = append(columns, m.viewRecent())
	columns = append(columns, m.viewTopStat(m.topSongs))
	columns = append(columns, m.viewTopStat(m.topArtists))
	columns = append(columns, m.viewTopStat(m.topGenres))

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func (m *Model) viewNotPlaying() string {
	rows := make([][]string, 0, 2)
	for i := 0; i < 2; i++ {
		rows = append(rows, make([]string, 0, 2))
	}

	rows[0] = append(rows[0], m.viewRecent())
	rows[0] = append(rows[0], m.viewTopStat(m.topSongs))
	rows[1] = append(rows[1], m.viewTopStat(m.topArtists))
	rows[1] = append(rows[1], m.viewTopStat(m.topGenres))

	renderedRows := make([]string, 0, 2)
	for _, row := range rows {
		renderedRows = append(renderedRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	view := lipgloss.JoinVertical(lipgloss.Left, renderedRows...)

	return sAll.Render(view)
}

func (m *Model) viewRecent() string {
	items := make([]string, 0, len(m.history))
	for i, track := range m.history {
		number := sStatEnum.Render(fmt.Sprintf("%d.", i+1))
		body := sStatBody.Render(track)
		items = append(items, lipgloss.JoinHorizontal(lipgloss.Top, number, body))
	}
	l := lipgloss.JoinVertical(lipgloss.Left, items...)
	title := sStatTitle.Render("Recently Played")

	return sStat.Render(lipgloss.JoinVertical(lipgloss.Left, title, l))
}

func (m *Model) viewTopStat(topStat topStat) string {
	items := make([]string, 0, len(topStat.entries))
	for i, stat := range topStat.entries {
		number := sStatEnum.Render(fmt.Sprintf("%d.", i+1))
		body := sStatBody.Render(stat.name)
		amount := sStatAmount.Render(fmt.Sprintf("%d", stat.amount))
		items = append(items, lipgloss.JoinHorizontal(lipgloss.Top, number, body, amount))
	}
	l := lipgloss.JoinVertical(lipgloss.Left, items...)

	title := sStatTitle.Render(topStat.title)

	return sStat.Render(lipgloss.JoinVertical(lipgloss.Left, title, l))
}
