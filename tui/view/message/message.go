// Package message provides the functions to draw all the cammie messages on a TUI
package message

import (
	"context"
	"hash/fnv"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
	"github.com/zeusWPI/scc/tui/view"
)

// Model represents the model for the message view
type Model struct {
	repo      repository.Message
	messages  []message
	blacklist []string

	lastMessageID int
	width         int
	height        int
}

// Interface compliance
var _ view.View = (*Model)(nil)

// Msg represents the message to update the message view
type Msg struct {
	lastMessageID int
	messages      []message
}

type message struct {
	sender  string
	message string
	color   string
	date    time.Time
}

func NewModel(repo repository.Repository) view.View {
	return &Model{
		repo:          *repo.NewMessage(),
		messages:      nil,
		blacklist:     config.GetDefaultStringSlice("cammie.blacklist", []string{}),
		lastMessageID: -1,
		width:         0,
		height:        0,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Name() string {
	return "Cammie Messages"
}

func (m *Model) Update(msg tea.Msg) (view.View, tea.Cmd) {
	switch msg := msg.(type) {
	case view.MsgSize:
		// Size update!
		// Check if it's relevant for this view
		if entry, ok := msg.Sizes[m.Name()]; ok {
			// Update all dependent styles
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
			Interval: config.GetDefaultInt("tui.view.message.interval_s", 1),
		},
	}
}

func updateMessages(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)
	lastMessageID := m.lastMessageID

	messagesDB, err := m.repo.GetSinceID(ctx, lastMessageID)
	if err != nil {
		return nil, err
	}

	if len(messagesDB) == 0 {
		return nil, nil
	}

	messages := make([]message, 0, len(messagesDB))
	lastID := m.lastMessageID
	for _, msg := range messagesDB {
		if slices.Contains(m.blacklist, msg.Name) {
			continue
		}

		if msg.ID > lastID {
			lastID = msg.ID
		}

		messages = append(messages, message{
			sender:  msg.Name,
			message: msg.Message,
			color:   hashColor(msg.Name),
			date:    msg.CreatedAt,
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
