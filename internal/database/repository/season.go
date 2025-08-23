package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/sqlc"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Season struct {
	repo Repository
}

func (r *Repository) NewSeason() *Season {
	return &Season{
		repo: *r,
	}
}

// TODO: Check still used
func (s *Season) GetAll(ctx context.Context) ([]*model.Season, error) {
	seasons, err := s.repo.queries(ctx).SeasonGetAll(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("season get all %w", err)
		}
		return nil, nil
	}

	return utils.SliceMap(seasons, model.SeasonModel), nil
}

func (s *Season) GetCurrent(ctx context.Context) (*model.Season, error) {
	season, err := s.repo.queries(ctx).SeasonGetCurrent(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("season get current %w", err)
		}
		return nil, nil
	}

	return model.SeasonModel(season), nil
}

func (s *Season) Create(ctx context.Context, season *model.Season) error {
	id, err := s.repo.queries(ctx).SeasonCreate(ctx, sqlc.SeasonCreateParams{
		Name:    season.Name,
		Start:   pgtype.Timestamp{Time: time.Time(season.Start), Valid: !time.Time(season.Start).IsZero()},
		End:     pgtype.Timestamp{Time: time.Time(season.End), Valid: !time.Time(season.End).IsZero()},
		Current: season.Current,
	})
	if err != nil {
		return fmt.Errorf("create season %+v | %w", *season, err)
	}

	season.ID = int(id)

	return nil
}

func (s *Season) Update(ctx context.Context, season model.Season) error {
	if err := s.repo.queries(ctx).SeasonUpdate(ctx, sqlc.SeasonUpdateParams{
		ID:      int32(season.ID),
		Name:    season.Name,
		Start:   pgtype.Timestamp{Time: time.Time(season.Start), Valid: !time.Time(season.Start).IsZero()},
		End:     pgtype.Timestamp{Time: time.Time(season.End), Valid: !time.Time(season.End).IsZero()},
		Current: season.Current,
	}); err != nil {
		return fmt.Errorf("update season %+v | %w", season, err)
	}

	return nil
}

func (s *Season) DeleteAll(ctx context.Context) error {
	if err := s.repo.queries(ctx).SeasonDeleteAll(ctx); err != nil {
		return fmt.Errorf("delete all seasons %w", err)
	}

	return nil
}
