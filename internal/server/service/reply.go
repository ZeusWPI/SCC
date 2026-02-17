package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/internal/server/dto"
	"go.uber.org/zap"
)

type Reply struct {
	message repository.Message
	reply   repository.Reply
}

func (s *Service) NewReply() *Reply {
	return &Reply{
		message: *s.repo.NewMessage(),
		reply:   *s.repo.NewReply(),
	}
}

func (r *Reply) Create(ctx context.Context, replySave dto.ReplySave) (dto.Reply, error) {
	reply := replySave.ToModel()

	msg, err := r.message.Get(ctx, reply.MessageID)
	if err != nil {
		zap.S().Error(err)
		return dto.Reply{}, fiber.ErrInternalServerError
	}
	if msg == nil {
		return dto.Reply{}, fiber.ErrNotFound
	}

	if err := r.reply.Create(ctx, reply); err != nil {
		zap.S().Error(err)
		return dto.Reply{}, fiber.ErrInternalServerError
	}

	return dto.ReplyDTO(reply), nil
}
