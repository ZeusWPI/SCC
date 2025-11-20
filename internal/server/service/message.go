package service

import (
	"context"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/buzzer"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/server/dto"
	"github.com/zeusWPI/scc/pkg/config"
	"go.uber.org/zap"
)

type Message struct {
	message repository.Message

	buzzer    buzzer.Client
	blacklist []string
}

func (s *Service) NewMessage() *Message {
	return &Message{
		message:   *s.repo.NewMessage(),
		buzzer:    *buzzer.New(),
		blacklist: config.GetDefaultStringSlice("cammie.blacklist", []string{}),
	}
}

func (m *Message) Create(ctx context.Context, msgSave dto.Message) (dto.Message, error) {
	msg := msgSave.ToModel()
	if err := m.message.Create(ctx, msg); err != nil {
		zap.S().Error(err)
		return dto.Message{}, fiber.ErrInternalServerError
	}

	if !slices.Contains(m.blacklist, msg.Name) {
		m.buzzer.Play()
	}

	return dto.MessageDTO(msg), nil
}
