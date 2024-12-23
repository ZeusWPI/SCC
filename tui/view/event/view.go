package event

import (
	"bytes"
	"image"

	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
)

func (m *Model) viewToday() string {
	// Render image
	poster := ""
	if m.today.Poster != nil {
		i, _, err := image.Decode(bytes.NewReader(m.today.Poster))
		if err == nil {
			poster = view.ImagetoString(wTodayPoster, i)
		}
	}

	name := sTodayText.Render(m.today.Name)
	date := sTodayDate.Render("🕙 " + m.today.Date.Format("15:04"))
	location := sTodayeLoc.Render("📍 " + m.today.Location)

	event := lipgloss.JoinVertical(lipgloss.Left, name, date, location)
	event = sToday.Render(event)

	if lipgloss.Height(poster) > lipgloss.Height(event) {
		event = sTodayPoster.Height(lipgloss.Height(poster)).Render(event)
	} else {
		poster = sTodayPoster.Height(lipgloss.Height(event)).Render(poster)
	}

	view := lipgloss.JoinHorizontal(lipgloss.Top, poster, event)

	return sTodayAll.Render(view)
}

func (m *Model) viewOverview() string {
	// Poster if present
	poster := ""
	if len(m.upcoming) > 0 && m.upcoming[0].Poster != nil {
		i, _, err := image.Decode(bytes.NewReader(m.upcoming[0].Poster))
		if err == nil {
			poster = view.ImagetoString(wOvPoster, i)
		}
	}

	// Overview
	events := m.viewGetEventOverview()

	if lipgloss.Height(poster) > lipgloss.Height(events) {
		events = sOv.Height(lipgloss.Height(poster)).Render(events)
	} else {
		poster = sOvPoster.Height(lipgloss.Height(events)).Render(poster)
	}

	// Combine image and overview
	view := lipgloss.JoinHorizontal(lipgloss.Top, events, poster)

	return sOvAll.Render(view)
}

func (m *Model) viewGetEventOverview() string {
	events := make([]string, 0, len(m.passed)+len(m.upcoming)+1)

	title := sOvTitle.Render("Events")
	events = append(events, title)

	// Passed
	for _, event := range m.passed {
		date := sOvPassedDate.Render(event.Date.Format("02/01"))
		name := sOvPassedText.Render(event.Name)
		text := lipgloss.JoinHorizontal(lipgloss.Top, date, name)

		events = append(events, text)
	}

	if len(m.upcoming) > 0 {
		// Next
		date := sOvNextDate.Render(m.upcoming[0].Date.Format("02/01"))
		name := sOvNextText.Render(m.upcoming[0].Name)
		location := sOvNextLoc.Render("📍 " + m.upcoming[0].Location)

		text := lipgloss.JoinVertical(lipgloss.Left, name, location)
		text = lipgloss.JoinHorizontal(lipgloss.Top, date, text)

		events = append(events, text)
	}

	// Upcoming
	for i := 1; i < len(m.upcoming); i++ {
		date := sOvUpcomingDate.Render(m.upcoming[i].Date.Format("02/01"))
		name := sOvUpcomingText.Render(m.upcoming[i].Name)
		text := name
		if i < 3 {
			location := sOvNextLoc.Render("📍 " + m.upcoming[i].Location)
			text = lipgloss.JoinVertical(lipgloss.Left, name, location)
		}

		text = lipgloss.JoinHorizontal(lipgloss.Top, date, text)

		events = append(events, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, events...)

	return sOv.Render(view)
}
