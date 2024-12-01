// Package message provides the functions to draw all the cammie messages on a TUI
package message

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/ui/view"
	"go.uber.org/zap"
)

// Model represents the model for the message view
type Model struct {
	db            *db.DB
	lastMessageID int64
	messages      []string
}

// Msg represents the message to update the message view
type Msg struct {
	lastMessageID int64
	messages      []string
}

var messageColor = []string{
	"#800000", "#008000", "#808000", "#000080", "#800080", "#008080", "#c0c0c0",
	"#ff0000", "#00ff00", "#ffff00", "#0000ff", "#ff00ff", "#00ffff", "#ffffff",
}

// NewModel creates a new message model view
func NewModel(db *db.DB) view.View {
	return &Model{db: db, lastMessageID: -1, messages: []string{}}
}

// Init initializes the message model view
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update updates the message model view
func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg:
		m.lastMessageID = msg.lastMessageID
		m.messages = append(m.messages, msg.messages...)

		return m, nil
	}

	return m, nil
}

// View returns the view for the message model
func (m *Model) View() string {
	// TODO: Limit the amount of messages shown
	// TODO: Wrap messages
	zap.S().Info("Viewing messages")
	l := list.New(m.messages).Enumerator(func(_ list.Items, _ int) string { return "" })
	return l.String()
}

// GetUpdateDatas returns all the update functions for the message model
func (m *Model) GetUpdateDatas() []view.UpdateData {
	return []view.UpdateData{
		{
			Name:     "cammie messages",
			View:     m,
			Update:   updateMessages,
			Interval: config.GetDefaultInt("tui.message.interval_s", 1),
		},
	}
}

func updateMessages(view view.View) (tea.Msg, error) {
	m := view.(*Model)
	lastMessageID := m.lastMessageID

	message, err := m.db.Queries.GetLastMessage(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return Msg{lastMessageID: lastMessageID, messages: []string{}}, err
	}

	if message.ID <= lastMessageID {
		return Msg{lastMessageID: lastMessageID, messages: []string{}}, nil
	}

	messages, err := m.db.Queries.GetMessageSinceID(context.Background(), lastMessageID)
	if err != nil {
		zap.S().Error("DB: Failed to get messages", err)
		return Msg{lastMessageID: lastMessageID, messages: []string{}}, err
	}

	formattedMessages := make([]string, 0, len(messages))
	for _, message := range messages {
		formattedMessages = append(formattedMessages, formatMessage(message))
	}

	return Msg{lastMessageID: message.ID, messages: formattedMessages}, nil
}

func hashColor(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	hash := h.Sum32()
	return messageColor[hash%uint32(len(messageColor))]
}

func formatMessage(msg sqlc.Message) string {
	dateStyle := lipgloss.NewStyle().Faint(true)
	date := dateStyle.Render(fmt.Sprintf("%s ", msg.CreatedAt.Format("02/01")))

	color := hashColor(msg.Name)
	colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))

	sender := fmt.Sprintf("%s %s ", colorStyle.Bold(true).Render(msg.Name), colorStyle.Render("|"))
	message := colorStyle.Render(msg.Message)

	return fmt.Sprintf("%s%s%s", date, sender, message)
}
