package dto

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

type Message struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	IP      string    `json:"ip"`
	Message string    `json:"message"`
	SendAt  time.Time `json:"send_at"`
	Replies []Reply   `json:"replies,omitzero"`
}

func MessageDTO(msg *model.Message) Message {
	return Message{
		ID:      msg.ID,
		Name:    msg.Name,
		IP:      msg.IP,
		Message: msg.Message,
		SendAt:  msg.CreatedAt,
		Replies: []Reply{},
	}
}

type MessageCluster struct {
	Messages []Message `json:"messages"`
}

type MessageDayGroup struct {
	DateKey   string           `json:"date_key"`
	DateLabel string           `json:"date_label"`
	Clusters  []MessageCluster `json:"clusters"`
}

type MessageSave struct {
	Name    string `json:"name" validate:"required"`
	IP      string `json:"ip" validate:"required"`
	Message string `json:"message" validate:"required"`
}

func (m *MessageSave) ToModel() *model.Message {
	return &model.Message{
		Name:    m.Name,
		IP:      m.IP,
		Message: m.Message,
	}
}
