// Package event provides all logic regarding the events of the website
package event

import (
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/zeusWPI/scc/internal/pkg/db"
	"github.com/zeusWPI/scc/internal/pkg/db/dto"
	"github.com/zeusWPI/scc/internal/pkg/db/sqlc"
	"github.com/zeusWPI/scc/pkg/config"
)

// Event represents a event instance
type Event struct {
	db            *db.DB
	website       string
	websitePoster string
}

// New creates a new event instance
func New(db *db.DB) *Event {
	return &Event{
		db:            db,
		website:       config.GetDefaultString("backend.event.website", "https://zeus.gent/events/"),
		websitePoster: config.GetDefaultString("backend.event.website_poster", "https://git.zeus.gent/ZeusWPI/visueel/raw/branch/master"),
	}
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

	eventsDBSQL, err := e.db.Queries.GetEventByAcademicYear(context.Background(), events[0].AcademicYear)
	if err != nil {
		return err
	}

	eventsDB := make([]*dto.Event, 0, len(eventsDBSQL))

	var wg sync.WaitGroup
	var errs []error
	for _, event := range eventsDBSQL {
		wg.Add(1)

		go func(event sqlc.Event) {
			defer wg.Done()

			ev := dto.EventDTO(event)
			eventsDB = append(eventsDB, ev)
			err := e.getPoster(ev)
			if err != nil {
				errs = append(errs, err)
			}
		}(event)
	}
	wg.Wait()

	// Check if there are any new events
	equal := true
	for _, event := range eventsDB {
		found := slices.ContainsFunc(events, func(ev dto.Event) bool { return ev.Equal(*event) })
		if !found {
			equal = false
			break
		}
	}

	if len(events) != len(eventsDB) {
		equal = false
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

	for _, event := range events {
		err = e.getPoster(&event)
		if err != nil {
			errs = append(errs, err)
			// Don't return / continue. We can still enter it without a poster
		}
		_, err = e.db.Queries.CreateEvent(context.Background(), event.CreateParams())
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
