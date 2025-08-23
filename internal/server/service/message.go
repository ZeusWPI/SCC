package service

import (
	"context"

	"github.com/zeusWPI/scc/internal/server/dto"
)

type Message struct{}

func (s *Service) NewMessage() *Message {
	return &Message{}
}

// TODO: fill in
func (m *Message) Create(ctx context.Context, message dto.Message) (dto.Message, error) {
	return dto.Message{}, nil
}
