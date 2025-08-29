package service

import (
	"context"

	"github.com/zeusWPI/scc/internal/server/dto"
)

type Message struct{}

func (s *Service) NewMessage() *Message {
	return &Message{}
}

func (m *Message) Create(_ context.Context, _ dto.Message) (dto.Message, error) {
	// TODO: fill in
	return dto.Message{}, nil
}
