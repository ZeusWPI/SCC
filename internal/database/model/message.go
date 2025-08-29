package model

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/sqlc"
)

type Message struct {
	ID        int
	Name      string
	IP        string
	Message   string
	CreatedAt time.Time
}

func MessageModel(m sqlc.Message) *Message {
	return &Message{
		ID:        int(m.ID),
		Name:      m.Name,
		IP:        m.Ip,
		Message:   m.Message,
		CreatedAt: m.CreatedAt.Time,
	}
}
