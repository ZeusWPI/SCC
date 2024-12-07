package message

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewAll() string {
	rows := make([]string, 0, len(m.messages))
	var previousDate time.Time

	for _, msg := range m.messages {
		time := sTime.Render(msg.date.Format("15:04") + " ")
		sender := sSender.Foreground(lipgloss.Color(msg.color)).Render(msg.sender + " | ")

		width := m.width - lipgloss.Width(time) - lipgloss.Width(sender)
		message := sMessage.Width(width).Foreground(lipgloss.Color(msg.color)).Render(msg.message)

		text := lipgloss.JoinHorizontal(lipgloss.Top, time, sender, message)

		// Add date if needed
		if previousDate.IsZero() || previousDate.YearDay() != msg.date.YearDay() {
			date := sDate.Render(" " + msg.date.Format("02/01") + " ")

			lineLength := (m.width - lipgloss.Width(date)) / 2
			left := sDate.Render(strings.Repeat("─", lineLength))
			right := sDate.Render(strings.Repeat("─", lineLength))

			date = lipgloss.JoinHorizontal(lipgloss.Top, left, date, right)
			text = lipgloss.JoinVertical(lipgloss.Left, date, text)
			previousDate = msg.date
		}

		rows = append(rows, text)
	}

	view := lipgloss.JoinVertical(lipgloss.Left, rows...)

	lines := strings.Split(view, "\n")
	height := min(m.height, len(lines))
	if len(lines) > 0 && len(lines) > height {
		view = strings.Join(lines[len(lines)-height:], "\n")
	}

	return view
}
