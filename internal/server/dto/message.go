package dto

import "github.com/zeusWPI/scc/internal/database/model"

type Message struct {
	Name    string `json:"name" validate:"required"`
	IP      string `json:"ip" validate:"required"`
	Message string `json:"message" validate:"required"`
}

func MessageDTO(msg *model.Message) Message {
	return Message{
		Name:    msg.Name,
		IP:      msg.IP,
		Message: msg.Message,
	}
}

func (m *Message) ToModel() *model.Message {
	return &model.Message{
		Name:    m.Name,
		IP:      m.IP,
		Message: m.Message,
	}
}
