package gamification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"slices"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zeusWPI/scc/pkg/utils"
	"github.com/zeusWPI/scc/tui/view"
)

func (g gamification) equal(g2 gamification) bool {
	return g.Name == g2.Name && g.Score == g2.Score && g.AvatarURL == g2.AvatarURL
}

func updateLeaderboard(ctx context.Context, view view.View) (tea.Msg, error) {
	m := view.(*Model)

	leaderboard, err := getLeaderboard(ctx, m.url)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(leaderboard, func(a, b gamification) int { return b.Score - a.Score })

	if len(leaderboard) != len(m.leaderboard) {
		return Msg{leaderboard: leaderboard}, nil
	}

	for idx, l := range leaderboard {
		if !m.leaderboard[idx].equal(l) {
			return Msg{leaderboard: leaderboard}, nil
		}
	}

	return nil, nil
}

func getLeaderboard(ctx context.Context, url string) ([]gamification, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    url + "/top4",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("get top 4 gamification %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var gams []gamification
	if err := json.NewDecoder(resp.Body).Decode(&gams); err != nil {
		return nil, fmt.Errorf("decode gamification response %w", err)
	}

	var errs []error

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := range gams {
		wg.Go(func() {
			if err := getAvatar(ctx, &gams[i]); err != nil {
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

	return gams, nil
}

func getAvatar(ctx context.Context, gam *gamification) error {
	if gam.AvatarURL == "" {
		return nil
	}

	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    gam.AvatarURL,
	})
	if err != nil {
		return fmt.Errorf("get avatar url %+v | %w", *gam, err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("decode gamification avatar %+v | %w", *gam, err)
	}

	gam.avatar = img

	return nil
}
