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

type Tap struct {
	repo Repository
}

func (r *Repository) NewTap() *Tap {
	return &Tap{
		repo: *r,
	}
}

func (t *Tap) GetLast(ctx context.Context) (*model.Tap, error) {
	tap, err := t.repo.queries(ctx).TapGetLast(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get last tap %w", err)
		}
		return nil, nil
	}

	return model.TapModel(tap), nil
}

func (t *Tap) GetCountByCategory(ctx context.Context) ([]*model.TapCount, error) {
	counts, err := t.repo.queries(ctx).TapGetCountByCategory(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get tap count by category %w", err)
		}
		return nil, nil
	}

	return utils.SliceMap(counts, model.TapCountModel), nil
}

func (t *Tap) Create(ctx context.Context, tap *model.Tap) error {
	id, err := t.repo.queries(ctx).TapCreate(ctx, sqlc.TapCreateParams{
		Name:           tap.Name,
		Category:       sqlc.TapCategory(tap.Category),
		OrderID:        int32(tap.OrderID),
		OrderCreatedAt: pgtype.Timestamptz{Time: tap.CreatedAt, Valid: !tap.CreatedAt.IsZero()},
	})
	if err != nil {
		return fmt.Errorf("create tap %+v | %w", *tap, err)
	}

	tap.ID = int(id)

	return nil
}
