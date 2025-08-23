package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/sqlc"
	"github.com/zeusWPI/scc/pkg/utils"
)

type Event struct {
	repo Repository
}

func (r *Repository) NewEvent() *Event {
	return &Event{
		repo: *r,
	}
}

func (e *Event) GetByAcademicYear(ctx context.Context, year string) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetAllByYear(ctx, year)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by academic year %s | %w", year, err)
		}
		return nil, nil
	}

	return utils.SliceMap(events, model.EventModel), nil
}

func (e *Event) GetByCurrentAcademicYear(ctx context.Context) ([]*model.Event, error) {
	events, err := e.repo.queries(ctx).EventGetAllByCurrentYear(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by current academic year %w", err)
		}
		return nil, nil
	}

	return utils.SliceMap(events, model.EventModel), nil
}

func (e *Event) Create(ctx context.Context, event *model.Event) error {
	id, err := e.repo.queries(ctx).EventCreate(ctx, sqlc.EventCreateParams{
		Name:         event.Name,
		Date:         pgtype.Timestamptz{Time: event.Date, Valid: !event.Date.IsZero()},
		AcademicYear: event.AcademicYear,
		Location:     event.Location,
		Poster:       event.Poster,
	})
	if err != nil {
		return fmt.Errorf("create event %+v | %w", *event, err)
	}

	event.ID = int(id)

	return nil
}

func (e *Event) DeleteAll(ctx context.Context) error {
	if err := e.repo.queries(ctx).EventDeleteAll(ctx); err != nil {
		return fmt.Errorf("delete all events %w", err)
	}

	return nil
}
