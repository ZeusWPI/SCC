package event

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

type event struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Start    time.Time `json:"start_time"`
	Poster   []byte
}

func (e event) equal(e2 event) bool {
	return e.ID == e2.ID && e.Name == e2.Name && e.Location == e2.Location && e.Start.Equal(e2.Start) && bytes.Equal(e.Poster, e2.Poster)
}

func updateEvents(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	events, err := getEvents(ctx, m.url)
	if err != nil {
		return nil, err
	}

	if len(events) != len(m.events) {
		return Msg{events: events}, nil
	}

	for _, ev := range events {
		if idx := slices.IndexFunc(m.events, func(e event) bool { return e.equal(ev) }); idx == -1 {
			return Msg{events: events}, nil
		}
	}

	return nil, nil
}

func getEvents(ctx context.Context, url string) ([]event, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/event", url),
	})
	if err != nil {
		return nil, fmt.Errorf("get events %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code for get events %s", resp.Status)
	}

	var events []event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decode event api %w", err)
	}

	var errs []error

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, event := range events {
		wg.Go(func() {
			if err := getPoster(ctx, url, &event); err != nil {
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

	return events, nil
}

func getPoster(ctx context.Context, url string, event *event) error {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/event/poster/%d?original=true&scc=true", url, event.ID),
	})
	if err != nil {
		return fmt.Errorf("get poster %+v | %w", *event, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			// Event doesn't have a poster
			return nil
		}
		return fmt.Errorf("bad response code for event poster %s | %+v", resp.Status, *event)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read poster bytes %+v | %w", *event, err)
	}

	event.Poster = bytes

	return nil
}
