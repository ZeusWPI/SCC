package event

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

const (
	amountOfPassed = 1
	amountOfFuture = 2
)

func (m *Model) viewToday() string {
	today, ok := utils.SliceFind(m.events, func(e event) bool { return utils.SameDay(e.Start, time.Now()) })
	if !ok {
		return ""
	}

	// Render image
	poster := ""
	if today.poster != nil {
		poster = view.ImageToString(today.poster, 0, 30)
	}

	name := sTodayText.Render(today.Name)
	date := sTodayDate.Render("🕙 " + today.Start.Format("15:04"))
	location := sTodayeLoc.Render("📍 " + today.Location)

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
	passed := utils.SliceGet(utils.SliceFilter(m.events, func(e event) bool { return e.Start.Before(time.Now()) }), amountOfPassed)
	upcoming := utils.SliceGet(utils.SliceFilter(m.events, func(e event) bool { return e.Start.After(time.Now()) }), amountOfFuture)

	// Poster if present
	poster := ""
	if len(upcoming) > 0 && upcoming[0].poster != nil {
		poster = view.ImageToString(upcoming[0].poster, 0, 30)
	}

	// Overview
	events := m.viewGetEventOverview(passed, upcoming)

	if lipgloss.Height(poster) > lipgloss.Height(events) {
		events = sOv.Height(lipgloss.Height(poster)).Render(events)
	} else {
		poster = sOvPoster.Height(lipgloss.Height(events)).Render(poster)
	}

	// Combine image and overview
	view := lipgloss.JoinHorizontal(lipgloss.Top, events, poster)

	return sOvAll.Render(view)
}

func (m *Model) viewGetEventOverview(passed, upcoming []event) string {
	events := make([]string, 0, len(passed)+len(upcoming)+1)

	title := sOvTitle.Render("Events")
	events = append(events, title)

	// Passed
	for _, event := range passed {
		date := sOvPassedDate.Render(event.Start.Format("02/01"))
		name := sOvPassedText.Render(event.Name)
		text := lipgloss.JoinHorizontal(lipgloss.Top, date, name)

		events = append(events, text)
	}

	if len(upcoming) > 0 {
		// Next
		date := sOvNextDate.Render(upcoming[0].Start.Format("02/01"))
		name := sOvNextText.Render(upcoming[0].Name)
		location := sOvNextLoc.Render("📍 " + upcoming[0].Location)

		text := lipgloss.JoinVertical(lipgloss.Left, name, location)
		text = lipgloss.JoinHorizontal(lipgloss.Top, date, text)

		events = append(events, text)
	}

	// Upcoming
	for i := 1; i < len(upcoming); i++ {
		date := sOvUpcomingDate.Render(upcoming[i].Start.Format("02/01"))
		name := sOvUpcomingText.Render(upcoming[i].Name)
		text := name
		if i < 3 {
			location := sOvNextLoc.Render("📍 " + upcoming[i].Location)
			text = lipgloss.JoinVertical(lipgloss.Left, name, location)
		}

		text = lipgloss.JoinHorizontal(lipgloss.Top, date, text)

		events = append(events, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, events...)

	return sOv.Render(view)
}
