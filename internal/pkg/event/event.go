// Package event provides all logic regarding the events of the website
package event

import (
	"context"
	"errors"
	"slices"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/pkg/config"
)

// Event represents a event instance
type Event struct {
	db  *db.DB
	api string
}

// New creates a new event instance
func New(db *db.DB) *Event {
	api := config.GetDefaultString("event.api", "https://zeus.gent/events")

	return &Event{db: db, api: api}
}

// Update gets all events from the website of this academic year
func (e *Event) Update() error {
	events, err := e.getEvents()
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	eventsDB, err := e.db.Queries.GetEventByAcademicYear(context.Background(), events[0].AcademicYear)
	if err != nil {
		return err
	}

	equal := false
	if len(events) == len(eventsDB) {
		for _, event := range eventsDB {
			found := slices.ContainsFunc(events, func(ev dto.Event) bool { return ev.Equal(*dto.EventDTO(event)) })
			if !found {
				break
			}
		}
		equal = true
	}

	// Both are equal, nothing to be done
	if equal {
		return nil
	}

	// They differ, remove the old ones and insert the new once
	err = e.db.Queries.DeleteEventByAcademicYear(context.Background(), events[0].AcademicYear)
	if err != nil {
		return err
	}
	var errs []error
	for _, event := range events {
		_, err = e.db.Queries.CreateEvent(context.Background(), event.CreateParams())
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
