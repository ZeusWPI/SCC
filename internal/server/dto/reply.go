package dto

import (
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
)

type Reply struct {
	ID      int       `json:"id"`
	Name    string    `json:"name,omitzero"`
	Message string    `json:"message"`
	SendAt  time.Time `json:"send_at"`
}

func ReplyDTO(r *model.Reply) Reply {
	return Reply{
		ID:      r.ID,
		Name:    r.Name,
		Message: r.Message,
		SendAt:  r.CreatedAt,
	}
}

type ReplySave struct {
	MessageID int    `json:"message_id" validate:"required"`
	Name      string `json:"name"`
	Message   string `json:"message" validate:"required,min=1"`
}

func (r ReplySave) ToModel() *model.Reply {
	return &model.Reply{
		MessageID: r.MessageID,
		Name:      r.Name,
		Message:   r.Message,
	}
}
