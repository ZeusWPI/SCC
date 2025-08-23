package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/sqlc"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Message struct {
	repo Repository
}

func (r *Repository) NewMessage() *Message {
	return &Message{
		repo: *r,
	}
}

func (m *Message) GetSinceID(ctx context.Context, id int) ([]*model.Message, error) {
	messages, err := m.repo.queries(ctx).MessageGetSinceID(ctx, int32(id))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get messages since id %d | %w", id, err)
		}
		return nil, nil
	}

	return utils.SliceMap(messages, model.MessageModel), nil
}

func (m *Message) Create(ctx context.Context, message *model.Message) error {
	id, err := m.repo.queries(ctx).MessageCreate(ctx, sqlc.MessageCreateParams{
		Name:    message.Name,
		Ip:      message.IP,
		Message: message.Message,
	})
	if err != nil {
		return fmt.Errorf("create message %w", err)
	}

	message.ID = int(id)

	return nil
}
