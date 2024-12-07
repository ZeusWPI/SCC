// Package message provides the functions to draw all the cammie messages on a TUI
package message

import (
	"context"
	"database/sql"
	"hash/fnv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/ui/view"
)

// Model represents the model for the message view
type Model struct {
	width         int
	height        int
	db            *db.DB
	lastMessageID int64
	messages      []message
}

type message struct {
	sender  string
	message string
	color   string
	date    time.Time
}

// Msg represents the message to update the message view
type Msg struct {
	lastMessageID int64
	messages      []message
}

// NewModel creates a new message model view
func NewModel(db *db.DB) view.View {
	return &Model{db: db, lastMessageID: -1, messages: []message{}}
}

// Init initializes the message model view
func (m *Model) Init() tea.Cmd {
	return nil
}

// Name returns the name of the view
func (m *Model) Name() string {
	return "Cammie Messages"
}

// Update updates the message model view
func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		entry, ok := msg.Sizes[m.Name()]
		if ok {
			m.width = entry.Width
			m.height = entry.Height
		}

		return m, nil
	case Msg:
		m.lastMessageID = msg.lastMessageID
		m.messages = append(m.messages, msg.messages...)

		return m, nil
	}

	return m, nil
}

// View returns the view for the message model
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	return m.viewAll()
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

	messagesDB, err := m.db.Queries.GetMessageSinceID(context.Background(), lastMessageID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return nil, err
	}

	if len(messagesDB) == 0 {
		return nil, nil
	}

	messages := make([]message, 0, len(messagesDB))
	lastID := m.lastMessageID
	for _, m := range messagesDB {
		if m.ID > lastID {
			lastID = m.ID
		}

		messages = append(messages, message{
			sender:  m.Name,
			message: m.Message,
			color:   hashColor(m.Name),
			date:    m.CreatedAt,
		})
	}

	return Msg{lastMessageID: lastID, messages: messages}, nil
}

func hashColor(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	hash := h.Sum32()
	return colors[hash%uint32(len(colors))]
}
