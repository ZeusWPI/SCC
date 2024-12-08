package dto

import (
	"time"

	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
)

// Message is the DTO for the message
type Message struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	IP        string    `json:"ip" validate:"required"`
	Message   string    `json:"message" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

// MessageDTO converts a sqlc.Message to a Message
func MessageDTO(message sqlc.Message) *Message {
	return &Message{
		ID:        message.ID,
		Name:      message.Name,
		IP:        message.Ip,
		Message:   message.Message,
		CreatedAt: message.CreatedAt.Time,
	}
}

// CreateParams converts a Message to sqlc.CreateMessageParams
func (m *Message) CreateParams() sqlc.CreateMessageParams {
	return sqlc.CreateMessageParams{
		Name:    m.Name,
		Ip:      m.IP,
		Message: m.Message,
	}
}
