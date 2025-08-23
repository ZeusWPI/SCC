package gamification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/zeusWPI/scc/internal/database/model"
	"github.com/zeusWPI/scc/pkg/utils"
)

type gamificationAPI struct {
	Name      string `json:"github_name"`
	Score     int    `json:"score"`
	AvatarURL string `json:"avatar_url"`
	Avatar    []byte
}

func (g gamificationAPI) toModel() model.Gamification {
	return model.Gamification{
		Name:   g.Name,
		Score:  g.Score,
		Avatar: g.Avatar,
	}
}

func (g *Gamification) getAvatar(ctx context.Context, gam *gamificationAPI) error {
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

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read avatar bytes %+v | %w", *gam, err)
	}

	gam.Avatar = bytes

	return nil
}

func (g *Gamification) getLeaderboard(ctx context.Context) ([]model.Gamification, error) {
	resp, err := utils.DoRequest(ctx, utils.DoRequestValues{
		Method: "GET",
		URL:    fmt.Sprintf("%s/top4", g.url),
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

	var gams []gamificationAPI
	if err := json.NewDecoder(resp.Body).Decode(&gams); err != nil {
		return nil, fmt.Errorf("decode gamification response %w", err)
	}

	var errs []error

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, gam := range gams {
		wg.Go(func() {
			if err := g.getAvatar(ctx, &gam); err != nil {
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

	return utils.SliceMap(gams, func(g gamificationAPI) model.Gamification { return g.toModel() }), nil
}
