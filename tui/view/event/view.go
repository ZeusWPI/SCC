package event

import (
	"bytes"
	"image"

	"github.com/charmbracelet/lipgloss"
	"github.com/zeusWPI/scc/tui/view"
)

func (m *Model) viewToday() string {
	// Render image
	im := ""
	if m.today.Poster != nil {
		i, _, err := image.Decode(bytes.NewReader(m.today.Poster))
		if err == nil {
			im = view.ImagetoString(widthImage, i)
		}
	}

	// Render text
	warningTop := sTodayWarning.MarginBottom(mTodayWarning).Render("🥳 Event Today 🥳")
	warningBottom := sTodayWarning.MarginTop(mTodayWarning).Render("🥳 Event Today 🥳")

	name := sTodayName.Render(m.today.Name)
	time := sTodayTime.Render("🕙 " + m.today.Date.Format("15:04"))
	location := sTodayPlace.Render("📍 " + m.today.Location)

	text := lipgloss.JoinVertical(lipgloss.Left, warningTop, name, time, location, warningBottom)

	// Resize so it's centered
	if lipgloss.Height(im) > lipgloss.Height(text) {
		sToday = sToday.Height(lipgloss.Height(im))
	}
	text = sToday.Render(text)

	return lipgloss.JoinHorizontal(lipgloss.Top, im, text)
}

func (m *Model) viewNormal() string {
	// Poster if present
	im := ""
	if len(m.upcoming) > 0 && m.upcoming[0].Poster != nil {
		i, _, err := image.Decode(bytes.NewReader(m.upcoming[0].Poster))
		if err == nil {
			im = view.ImagetoString(widthOverviewImage, i)
		}
	}

	// Overview
	events := m.viewGetEvents()

	// Filthy hack to avoid the last event being centered by the cammie screen
	events = append(events, "\n")

	// Render events overview
	overview := lipgloss.JoinVertical(lipgloss.Left, events...)
	overview = sOverview.Render(overview)

	title := sOverviewTitle.Render("Events")
	overview = lipgloss.JoinVertical(lipgloss.Left, title, overview)

	// Center the overview
	if lipgloss.Height(im) > lipgloss.Height(overview) {
		overview = sOverviewTotal.Height(lipgloss.Height(im)).Render(overview)
	}

	// Combine image and overview
	view := lipgloss.JoinHorizontal(lipgloss.Top, overview, im)

	return view
}

func (m *Model) viewGetEvents() []string {
	events := make([]string, 0, len(m.passed)+len(m.upcoming))

	// Passed
	for _, event := range m.passed {
		time := sPassedTime.Render(event.Date.Format("02/01") + "\t")
		name := sPassedName.Render(event.Name)
		text := lipgloss.JoinHorizontal(lipgloss.Top, time, name)

		events = append(events, text)
	}

	if len(m.upcoming) == 0 {
		return events
	}

	// Next
	name := sNextName.Render(m.upcoming[0].Name)
	time := sNextTime.Render(m.upcoming[0].Date.Format("02/01") + "\t")
	location := sNextPlace.Render("📍 " + m.upcoming[0].Location)

	text := lipgloss.JoinVertical(lipgloss.Left, name, location)
	text = lipgloss.JoinHorizontal(lipgloss.Top, time, text)

	events = append(events, text)

	// Upcoming
	for i := 1; i < len(m.upcoming); i++ {
		time := sUpcomingTime.Render(m.upcoming[i].Date.Format("02/01") + "\t")
		name := sUpcomingName.Render(m.upcoming[i].Name)
		text := name
		if i < 3 {
			location := sUpcomingPlace.Render("📍 " + m.upcoming[i].Location)
			text = lipgloss.JoinVertical(lipgloss.Left, name, location)
		}

		text = lipgloss.JoinHorizontal(lipgloss.Top, time, text)

		events = append(events, text)
	}

	return events
}
