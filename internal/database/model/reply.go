package model

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/sqlc"
)

type Reply struct {
	ID        int
	MessageID int
	Name      string
	Message   string
	CreatedAt time.Time
}

func ReplyModel(r sqlc.Reply) *Reply {
	name := ""
	if r.Name.Valid {
		name = r.Name.String
	}

	return &Reply{
		ID:        int(r.ID),
		MessageID: int(r.MessageID),
		Name:      name,
		Message:   r.Message,
		CreatedAt: r.CreatedAt.Time,
	}
}
