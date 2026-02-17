package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/internal/database/sqlc"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Reply struct {
	repo Repository
}

func (r *Repository) NewReply() *Reply {
	return &Reply{
		repo: *r,
	}
}

func (r *Reply) GetSinceMessageID(ctx context.Context, messageID int) ([]*model.Reply, error) {
	replies, err := r.repo.queries(ctx).ReplyGetSinceMessageID(ctx, int32(messageID))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get replies since message id %d | %w", messageID, err)
		}
		return nil, nil
	}

	return utils.SliceMap(replies, model.ReplyModel), nil
}

func (r *Reply) Create(ctx context.Context, reply *model.Reply) error {
	id, err := r.repo.queries(ctx).ReplyCreate(ctx, sqlc.ReplyCreateParams{
		MessageID: int32(reply.MessageID),
		Name:      pgtype.Text{String: reply.Name, Valid: reply.Name != ""},
		Message:   reply.Message,
	})
	if err != nil {
		return fmt.Errorf("create reply %+v | %w", *reply, err)
	}

	reply.ID = int(id)

	return nil
}
