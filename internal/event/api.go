package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/utils"
)

type eventAPI struct {
	ID        int       `json:"id"`
	Name      string    `json:"Name"`
	Location  string    `json:"location"`
	Start     time.Time `json:"start_time"`
	End       time.Time `json:"end_time"`
	YearStart int       `json:"year_start"`
	YearEnd   int       `json:"year_end"`
	Poster    []byte
}

func (e eventAPI) toModel() model.Event {
	return model.Event{
		Name:         e.Name,
		Date:         e.Start,
		AcademicYear: fmt.Sprintf("%d-%d", e.YearStart, e.YearEnd),
		Location:     e.Location,
		Poster:       e.Poster,
	}
}

func (e *Event) getPoster(ctx context.Context, event *eventAPI) error {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/event/poster/%d?original=true&scc=true", e.url, event.ID),
	})
	if err != nil {
		return fmt.Errorf("get poster %+v | %w", *event, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read poster bytes %+v | %w", *event, err)
	}

	event.Poster = bytes

	return nil
}

func (e *Event) getEvents(ctx context.Context) ([]model.Event, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/event", e.url),
	})
	if err != nil {
		return nil, fmt.Errorf("get events %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var events []eventAPI
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decode event api %w", err)
	}

	var errs []error

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, event := range events {
		wg.Go(func() {
			if err := e.getPoster(ctx, &event); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		})
	}

	wg.Wait()

	if errs != nil {
		return nil, errors.Join(errs...)
	}

	return utils.SliceMap(events, func(e eventAPI) model.Event { return e.toModel() }), nil
}
