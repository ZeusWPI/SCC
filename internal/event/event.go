// Package event keeps the current event database in sync
package event

import (
	"context"

	"github.com/zeusWPI/scc/internal/database/repository"
	"github.com/zeusWPI/scc/pkg/config"
)

type Event struct {
	repo  repository.Repository
	event repository.Event
	url   string
}

// TODO: Check if colly is removed

func New(repo repository.Repository) *Event {
	return &Event{
		repo:  repo,
		event: *repo.NewEvent(),
		url:   config.GetDefaultString("backend.event.url", "https://events.zeus.gent/api/v1"),
	}
}

// Update gets all events from the website of this academic year
func (e *Event) Update(ctx context.Context) error {
	events, err := e.getEvents(ctx)
	if err != nil {
		return err
	}

	return e.repo.WithRollback(ctx, func(ctx context.Context) error {
		if err := e.event.DeleteAll(ctx); err != nil {
			return err
		}

		for _, event := range events {
			if err := e.event.Create(ctx, &event); err != nil {
				return err
			}
		}

		return nil
	})
}
