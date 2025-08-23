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

type Gamification struct {
	repo Repository
}

func (r *Repository) NewGamification() *Gamification {
	return &Gamification{
		repo: *r,
	}
}

func (g *Gamification) GetAll(ctx context.Context) ([]*model.Gamification, error) {
	gams, err := g.repo.queries(ctx).GamificationGetAll(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get all gamification %w", err)
		}
		return nil, nil
	}

	return utils.SliceMap(gams, model.GamificationModel), nil
}

func (g *Gamification) Create(ctx context.Context, gam *model.Gamification) error {
	id, err := g.repo.queries(ctx).GamificationCreate(ctx, sqlc.GamificationCreateParams{
		Name:   gam.Name,
		Score:  int32(gam.Score),
		Avatar: gam.Avatar,
	})
	if err != nil {
		return fmt.Errorf("create gamification %+v | %w", *gam, err)
	}

	gam.ID = int(id)

	return nil
}

func (g *Gamification) DeleteAll(ctx context.Context) error {
	if err := g.repo.queries(ctx).GamificationDeleteAll(ctx); err != nil {
		return fmt.Errorf("delete all gamifications %w", err)
	}

	return nil
}
