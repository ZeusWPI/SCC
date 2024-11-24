package view

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"go.uber.org/zap"
)

// MessageModel represents the model for the message view
type MessageModel struct {
	db            *db.DB
	lastMessageID int64
	messages      []string
}

// MessageMsg represents the message to update the message view
type MessageMsg struct {
	lastMessageID int64
	messages      []string
}

var messageColor = []string{
	"#800000", "#008000", "#808000", "#000080", "#800080", "#008080", "#c0c0c0",
	"#ff0000", "#00ff00", "#ffff00", "#0000ff", "#ff00ff", "#00ffff", "#ffffff",
}

// NewMessageModel creates a new message model view
func NewMessageModel(db *db.DB) *MessageModel {
	return &MessageModel{db: db, lastMessageID: -1, messages: []string{}}
}

// Init initializes the message model view
func (c *MessageModel) Init() tea.Cmd {
	return updateMessages(c.db, c.lastMessageID)
}

// Update updates the message model view
func (c *MessageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MessageMsg:
		c.lastMessageID = msg.lastMessageID
		c.messages = append(c.messages, msg.messages...)

		return c, updateMessages(c.db, c.lastMessageID)
	}

	return c, nil
}

// View returns the view for the message model
func (c *MessageModel) View() string {
	// TODO: Limit the amount of messages shown
	// TODO: Wrap messages
	zap.S().Info("Viewing messages")
	l := list.New(c.messages).Enumerator(func(_ list.Items, _ int) string { return "" })
	return l.String()
}

func updateMessages(db *db.DB, lastMessageID int64) tea.Cmd {
	return tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
		message, err := db.Queries.GetLastMessage(context.Background())
		if err != nil {
			if err != sql.ErrNoRows {
				zap.S().Error("DB: Failed to get last message", err)
			}
			return MessageMsg{lastMessageID: lastMessageID, messages: []string{}}
		}

		if message.ID <= lastMessageID {
			return MessageMsg{lastMessageID: lastMessageID, messages: []string{}}
		}

		messages, err := db.Queries.GetMessageSinceID(context.Background(), lastMessageID)
		if err != nil {
			zap.S().Error("DB: Failed to get messages", err)
			return MessageMsg{lastMessageID: lastMessageID, messages: []string{}}
		}

		formattedMessages := make([]string, 0, len(messages))
		for _, message := range messages {
			formattedMessages = append(formattedMessages, formatMessage(message))
		}

		return MessageMsg{lastMessageID: message.ID, messages: formattedMessages}
	})
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
