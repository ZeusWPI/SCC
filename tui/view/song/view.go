package song

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/pkg/utils"
)

func (m *Model) viewPlaying() string {
	status := m.viewPlayingStatus()
	status = sStatus.Render(status)

	stats := m.viewPlayingStats()
	stats = sStat.Render(stats)

	lyrics := m.viewPlayingLyrics()
	lyrics = sLyric.Height(sAll.GetHeight() - lipgloss.Height(status) - lipgloss.Height(stats)).
		MaxHeight(sAll.GetHeight() - lipgloss.Height(status) - lipgloss.Height(stats)).
		Render(lyrics)

	view := lipgloss.JoinVertical(lipgloss.Left, status, lyrics, stats)

	return sAll.Render(view)
}

func (m *Model) viewPlayingStatus() string {
	// Stopwatch
	durationS := int(math.Round(float64(m.current.song.DurationMS) / 1000))
	stopwatch := fmt.Sprintf("\t%s / %02d:%02d", m.progress.stopwatch.View(), durationS/60, durationS%60)
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
	progress := m.progress.bar.View()
	progress = sStatusBar.Render(progress)

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

	columns = append(columns, m.viewStat(m.history, m.statAmountPlaying))
	columns = append(columns, m.viewStat(m.statsMonthly[0], m.statAmountPlaying))
	columns = append(columns, m.viewStat(m.statsMonthly[1], m.statAmountPlaying))
	columns = append(columns, m.viewStat(m.statsMonthly[2], m.statAmountPlaying))

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func (m *Model) viewNotPlaying() string {
	// Render stats
	rows := make([][]string, 0, 3)
	for range 3 {
		rows = append(rows, make([]string, 0, 2))
	}

	rows[0] = append(rows[0], m.viewStat(m.statsMonthly[0], m.statAmount, "Monthly"))
	rows[0] = append(rows[0], m.viewStat(m.stats[0], m.statAmount, "All Time"))
	rows[1] = append(rows[1], m.viewStat(m.statsMonthly[1], m.statAmount, "Monthly"))
	rows[1] = append(rows[1], m.viewStat(m.stats[1], m.statAmount, "All Time"))
	rows[2] = append(rows[2], m.viewStat(m.statsMonthly[2], m.statAmount, "Monthly"))
	rows[2] = append(rows[2], m.viewStat(m.stats[2], m.statAmount, "All Time"))

	renderedRows := make([]string, 0, 3)
	var title string
	for i, row := range rows {
		r := lipgloss.JoinHorizontal(lipgloss.Top, row...)
		title = sStatCategory.Render(sStatCategoryTitle.Render(m.stats[i].title)) // HACK: Make border same size as 2 stats next to each other
		renderedRows = append(renderedRows, lipgloss.JoinVertical(lipgloss.Left, title, r))
	}

	v := lipgloss.JoinVertical(lipgloss.Left, renderedRows...)

	// Render history
	items := make([]string, 0, len(m.history.entries))

	// Push it down
	for range lipgloss.Height(title) {
		items = append(items, "")
	}
	items = append(items, sStatTitle.Render(m.history.title))

	for i, entry := range utils.SliceGet(m.history.entries, m.statAmount*3+(4*2)) {
		enum := sStatEnum.Render(fmt.Sprintf("%d.", i+1))
		body := sStatEntry.Render(entry.name)
		amount := sStatAmount.Render(strconv.Itoa(entry.amount))
		items = append(items, lipgloss.JoinHorizontal(lipgloss.Top, enum, body, amount))
	}
	items = append(items, "") // HACK: Avoid the last item shifting to the right
	list := lipgloss.JoinVertical(lipgloss.Left, items...)
	// title := sStatTitle.Render(m.history.title)
	history := sStatHistory.Height(lipgloss.Height(v) - 1).MaxHeight(lipgloss.Height(v) - 1).Render(list) // - 1 to compensate for the hack newline at the end

	v = lipgloss.JoinHorizontal(lipgloss.Top, history, v)

	return sAll.Render(v)
}

func (m *Model) viewStat(stat stat, amount int, titleOpt ...string) string {
	entries := utils.SliceGet(stat.entries, amount)

	title := stat.title
	if len(titleOpt) > 0 {
		title = titleOpt[0]
	}

	items := make([]string, 0, len(entries))
	for i := range entries {
		enum := sStatEnum.Render(fmt.Sprintf("%d.", i+1))
		body := sStatEntry.Render(stat.entries[i].name)
		amount := sStatAmount.Render(strconv.Itoa(stat.entries[i].amount))

		items = append(items, lipgloss.JoinHorizontal(lipgloss.Top, enum, body, amount))
	}
	items = append(items, "") // HACK: Avoid the last item shifting to the right
	l := lipgloss.JoinVertical(lipgloss.Left, items...)

	t := sStatTitle.Render(title)

	return sStatOne.Render(lipgloss.JoinVertical(lipgloss.Left, t, l))
}
